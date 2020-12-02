/**
 * Copyright (c) 2019 TIBCO Software Inc.
 * All rights reserved.
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
 *
 * <p/>
 * File name: DataCryptoGrapher.java
 * Created on: 3/19/19
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: DataCryptoGrapher.java 3880 2020-04-16 22:35:51Z nimish $
 */


package com.tibco.tgdb.channel.impl;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.security.AlgorithmParameters;
import java.security.PublicKey;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.security.interfaces.DSAPublicKey;
import java.security.interfaces.ECPublicKey;
import java.security.spec.DSAParameterSpec;
import java.security.spec.ECParameterSpec;

import javax.crypto.Cipher;
import javax.crypto.interfaces.DHPublicKey;
import javax.crypto.spec.DHParameterSpec;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.impl.ProtocolDataOutputStream;

public class DataCryptoGrapher {

    long sessionId;
    X509Certificate remoteCert;
    PublicKey publicKey;
    AlgorithmParameters algparams;

    DataCryptoGrapher(long sessionId, byte[] certBuffer) throws TGException
    {
        //The certificate buffer contains the name, and certificate byte array
        try {
            this.sessionId = sessionId;
            ByteArrayInputStream dis = new ByteArrayInputStream(certBuffer);

            CertificateFactory cf = CertificateFactory.getInstance("X.509");
            this.remoteCert = (X509Certificate) cf.generateCertificate(dis);
            publicKey = this.remoteCert.getPublicKey();
            algparams = getAlgorithmParameters(publicKey);
            System.out.printf("Certificate Name :%s Algorithm : %s\n", "mod", publicKey.getFormat());

        }
        catch (Exception e) {
            throw new TGException(e);
        }

    }

    public byte[] encrypt(byte[] data) throws TGException
    {
        try {
            Cipher cipher = Cipher.getInstance(publicKey.getAlgorithm());
            cipher.init(Cipher.ENCRYPT_MODE, publicKey, algparams);
            return cipher.doFinal(data);
        }
        catch (Exception e) {
            throw new TGException(e);
        }

    }

    public byte[] decrypt(TGInputStream is) throws TGException
    {

        try {
            //byte[] buf = is.readBytes();
            ProtocolDataOutputStream out = new ProtocolDataOutputStream();

            long rand = is.readLong();
            long len = is.readLong();
            long cnt = len / 8;
            long rem = len % 8;
            for (long i=0; i<cnt; i++) {
                long v = is.readLong();
                long org = v ^ rand;
                out.writeLongAsBytes(org);
            }
            for (long i=0; i<rem; i++) {
                int b = is.readByte();
                out.writeByte(b);
            }

            return out.toByteArray();

            //return buf;
        }
        catch (IOException ioe) {
            throw new TGException(ioe);
        }

    }

    private AlgorithmParameters getAlgorithmParameters(PublicKey publicKey) throws Exception
    {
        if (publicKey == null) return null;

        if (publicKey instanceof DSAPublicKey) {
            AlgorithmParameters algparams = AlgorithmParameters.getInstance(publicKey.getAlgorithm());
            DSAPublicKey dsakey = (DSAPublicKey) publicKey;
            DSAParameterSpec dsaParams = (DSAParameterSpec) dsakey.getParams();
            algparams.init(dsaParams);
            return algparams;
        }

        if (publicKey instanceof ECPublicKey) {
            AlgorithmParameters algparams = AlgorithmParameters.getInstance(publicKey.getAlgorithm());
            ECPublicKey eckey = (ECPublicKey) publicKey;
            ECParameterSpec ecParams = (ECParameterSpec) eckey.getParams();
            algparams.init(ecParams);
            return algparams;
        }

        if (publicKey instanceof DHPublicKey) {
            AlgorithmParameters algparams = AlgorithmParameters.getInstance(publicKey.getAlgorithm());
            DHPublicKey dhkey = (DHPublicKey) publicKey;
            DHParameterSpec dhParams = (DHParameterSpec) dhkey.getParams();
            algparams.init(dhParams);
            return algparams;
        }

        return null;  //RSA doesn't have
    }
}
