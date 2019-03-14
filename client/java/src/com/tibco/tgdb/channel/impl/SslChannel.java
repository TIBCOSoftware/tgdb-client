
/**
 * Copyright 2018 TIBCO Software Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * <p/>
 * File name :SSLChannel
 * Created on: 12/16/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: SSLChannel.java 583 2016-03-15 02:02:39Z vchung $
 */

package com.tibco.tgdb.channel.impl;

import com.tibco.tgdb.channel.TGChannelUrl;
import com.tibco.tgdb.exception.TGChannelDisconnectedException;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGMessage;
import com.tibco.tgdb.pdu.impl.SessionForcefullyTerminated;
import com.tibco.tgdb.utils.ConfigName;
import com.tibco.tgdb.utils.SortedProperties;
import com.tibco.tgdb.utils.TGProperties;

import javax.net.ssl.*;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.security.*;
import java.security.cert.Certificate;
import java.security.cert.CertificateException;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.util.Collection;
import java.util.EnumMap;

import static com.tibco.tgdb.utils.HexUtils.formatHex;

public class SSLChannel extends TcpChannel {

    private final static String TLSV12PROTOCOL = "TLSv1.2";
    private final static String[] TLSProtocols = new String[] {TLSV12PROTOCOL};

    private SSLContext sslContext;
    private KeyManager keyManagers[];
    private TrustManager[] trustManagers;
    KeyStore keyStore;
    SSLSocketFactory sslSocketFactory;
    SSLSocket sslSocket;
    String[] supportedSuites;
    boolean verifyDbName;


    protected SSLChannel(TGChannelUrl linkUrl, TGProperties<String,String> props) throws TGException
    {
        super(linkUrl, props);
        initSSL();
    }

    private void initSSL() throws TGException
    {
        try {
            initKeyStore();
            initKeyManagers();
            initTrustManagers();
            sslContext = SSLContext.getInstance(TLSV12PROTOCOL);
            sslContext.init(keyManagers, trustManagers, new SecureRandom());
            sslSocketFactory = sslContext.getSocketFactory();
        }
        catch (Exception e) {
            throw TGException.buildException("initSSL failed", "TGDB-INITSSL-ERROR", e);
        }
    }

    private void initKeyManagers() throws Exception
    {
        KeyManagerFactory kmf = KeyManagerFactory.getInstance(KeyManagerFactory.getDefaultAlgorithm());
        String ksp = this.properties.getProperty(ConfigName.KeyStorePassword);
        kmf.init(keyStore, ksp==null?null:ksp.toCharArray());
        keyManagers = kmf.getKeyManagers();
    }

    private void initTrustManagers() throws Exception
    {
        verifyDbName = this.properties.getPropertyAsBoolean(ConfigName.TlsVerifyDatabaseName);
        if (verifyDbName) {
            trustManagers = new TrustManager[] { new DatabaseTrustManager(this.properties)};
            return;
        }
        TrustManagerFactory tmf  = TrustManagerFactory.getInstance(TrustManagerFactory.getDefaultAlgorithm());
        tmf.init(keyStore);
        trustManagers = tmf.getTrustManagers();
    }

    private void initKeyStore() throws Exception {
        keyStore = KeyStore.getInstance("JKS");
        //Load System certificate
        String sysTrustFile  = String.format("%s%slib%ssecurity%scacerts", System.getProperties().getProperty("java.home"),
                File.separator, File.separator, File.separator);

         keyStore.load(new FileInputStream(sysTrustFile), "changeit".toCharArray());

        //load the user defined Certificates.
        String trustedCerts = this.properties.getProperty(ConfigName.TlsTrustedCertificates);
        if (trustedCerts == null) return;
        String[] paths = trustedCerts.split(",");
        CertificateFactory cf = CertificateFactory.getInstance("X.509");
        for (String path : paths) {

            //keyStore.load(new FileInputStream(path), null);
            Collection<? extends Certificate> certs = cf.generateCertificates(new FileInputStream(path));
            for (Certificate c : certs) {
                X509Certificate cert = (X509Certificate) c;
                EnumMap<X509Name, String> emap = X509Name.parse(cert.getSubjectDN());
                keyStore.setCertificateEntry(emap.get(X509Name.CommonName), c);
            }
        }

      return ;
    }

    protected void createSocket() throws TGException
    {

        super.createSocket();
        try {
            sslSocket = (SSLSocket)sslSocketFactory.createSocket(this.socket, getHost(), getPort(), true);
            String[] suites = sslSocket.getEnabledCipherSuites();
            supportedSuites = TGCipherSuite.filterSuites(suites);
            sslSocket.setEnabledCipherSuites(supportedSuites);
            sslSocket.setEnabledProtocols(TLSProtocols);
            sslSocket.setUseClientMode(true);
        }
        catch (IOException ioe) {
            throw new TGChannelDisconnectedException(ioe);
        }
    }

    public void onConnect() throws TGException
    {
        try {
            TGMessage msg = tryRead(); //For ChannelDisconnected message
            if (msg instanceof SessionForcefullyTerminated) {
                throw new TGChannelDisconnectedException(((SessionForcefullyTerminated)msg).getKillString());
            }
            sslSocket.startHandshake();
            //Check HostVerifier
            this.setSocket(sslSocket);
            performHandshake(true);
            doAuthenticate();
        }
        catch (IOException ioe) {
            throw new TGChannelDisconnectedException(ioe);
        }
    }

    private static class DatabaseTrustManager implements X509TrustManager
    {
        TGProperties properties;
        String databaseName;
        DatabaseTrustManager(TGProperties properties) {
            this.properties = properties;
            databaseName = (String) properties.getProperty(ConfigName.ConnectionDatabaseName);
        }

        private String getCertificateFingerPrint(X509Certificate cert, boolean bDebug)
        {
            try {
                byte[] sig = cert.getEncoded();
                if (bDebug) {
                    String hexString = formatHex(sig);
                    System.out.printf("----- Encoded Certificate : \n%s\n", hexString);
                    System.out.println("------------------");
                }

                MessageDigest digest = MessageDigest.getInstance("SHA-256");
                byte[] encodedhash = digest.digest(sig);
                String hexString = formatHex(encodedhash);
                if (bDebug) System.out.println(hexString);
                return hexString;
            }
            catch (Exception e) {
                return null;
            }
        }
        public X509Certificate[] getAcceptedIssuers() {return null;}

        public void checkClientTrusted(X509Certificate[] certs, String authType) throws CertificateException {
            System.out.printf("%s:%s\n", certs, authType);
        }
        public void checkServerTrusted(X509Certificate[] certs, String authType) throws CertificateException
        {
            X509Certificate cert = certs[0];
            String fingerPrint = getCertificateFingerPrint(cert, false); //Future we can compare the fingerPrint Id.

            String dbname = databaseName + ".ce$t"; //This is fully qualified name
            String issuername = databaseName + ".1nit";
            Principal subject = certs[0].getSubjectDN();
            Principal issuer  = certs[0].getIssuerDN();
            EnumMap<X509Name, String> subjectComponent = X509Name.parse(subject);
            EnumMap<X509Name, String> issuerComponent = X509Name.parse(issuer);
            String certName = subjectComponent.get(X509Name.CommonName);
            String certIssuer = issuerComponent.get(X509Name.CommonName);

            if (dbname.equals(certName) && (issuername.equals(certIssuer))) return;
            throw new CertificateException(String.format("SSL Certificate Exception. Connecting to the Wrong database:Expecting :%s, Server Subject Name:%s", databaseName, certName));
        }
    }

    static enum X509Name
    {

        CommonName("cn"),
        Country("c"),
        EmailAddress("emailaddress"),
        OrganizationalUnit("ou"),
        Organization("o"),
        Locality("l"),
        State("st");


        private String key;

        X509Name(String key ) {
            this.key = key;
        }

        static X509Name getByKey(String key)
        {
            for (X509Name n : X509Name.values())
            {
                if (n.key.equalsIgnoreCase(key)) return n;
            }
            return null;
        }

        public static EnumMap<X509Name, String> parse(Principal subject)
        {
            EnumMap<X509Name, String> map = new EnumMap<X509Name, String>(X509Name.class);

            String subjName = subject.getName();
            String tokens[] = subjName.split(",");
            for (String token:tokens) {
                String nvpair[] = token.split("=");
                X509Name name = X509Name.getByKey(nvpair[0].trim());
                if (name != null) {
                    map.put(name, nvpair[1].trim());
                }
                else {
                    System.out.printf("Found key :%s\n", nvpair[0]);
                }
            }
            return map;
        }


    }




    public static void main(String[] args) throws Exception {
        Provider[] providers = Security.getProviders();

        Class<Provider> clz = (Class<Provider>) Class.forName("com.sun.net.ssl.internal.ssl.Provider");
        Boolean b = (Boolean) clz.getDeclaredMethod("isFIPS").invoke(null);

        Provider pp = (Provider) Class.forName("com.sun.net.ssl.internal.ssl.Provider").newInstance();

        //Provider pp = Security.getProvider("com.sun.net.ssl.internal.ssl.Provider");
        System.out.printf("Class: %s, isFips:%s\n",pp.getClass().getName(), pp.keySet());

        // How to use JKS
        // https://blog.codecentric.de/en/2013/01/how-to-use-self-signed-pem-client-certificates-in-java/

        LinkUrl url = (LinkUrl) LinkUrl.parse("tcp://scott@foo.bar.com:8700");
        TGProperties<String,String> props = new SortedProperties<>(url.getProperties());
        props.put(ConfigName.TlsTrustedCertificates.getAlias(),
                "/Users/suresh/tsi-root/ext/openssl/1.0.2g/bin/suresh-mbp.pem,/Users/suresh/tibco/ems/8.0/samples/certs/server_root.cert.pem");
        SSLChannel sslChannel = new SSLChannel(url, props);
        /*
        for (int i=0; i<providers.length; i++) {
            Provider p = providers[i];

            System.out.printf("%s:%s\n", p.getName(), p.getClass().getName());
            System.out.println("=============================");
//            for(Provider.Service ps : p.getServices()) {
//                System.out.printf("\t%s[desc=%s,algo=%s]\n",ps, ps.getType(), ps.getAlgorithm());
//            }
        }
        */
        /*
        SSLContext ctx = SSLContext.getInstance("TLSv1.2");
        ctx.init(null, null, null);
        String[] suites = ctx.createSSLEngine().getSupportedCipherSuites();
        for (String suite : suites) {
            System.out.println(suite);
        }
        */

    }

}
