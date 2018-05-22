package com.tibco.tgdb.channel.impl;

import com.tibco.tgdb.channel.TGChannelUrl;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.utils.ConfigName;
import com.tibco.tgdb.utils.TGEnvironment;

import java.io.Console;
import java.io.FileDescriptor;
import java.io.FileInputStream;
import java.io.Reader;
import java.util.*;

/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
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
 * File name :LinkUrl
 * Created on: 12/16/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: LinkUrl.java 2203 2018-04-04 01:58:36Z ssubrama $
 */
public class LinkUrl implements TGChannelUrl
{

    //From http://www.networksorcery.com/enp/protocol/ip/ports08000.htm - looks like no well-established software uses this port.
    private final static String gDefaultHost = "localhost";
    private final static int    gDefaultPort = 8222;
    private final static LinkUrl gDefaultUrl = new LinkUrl(Protocol.TCP, gDefaultHost, gDefaultPort);
    private final static LinkedList<TGChannelUrl> EMPTY_LINKEDLIST = new LinkedList<TGChannelUrl>();

    String user;
    Protocol protocol;
    String urlstr = null;
    String              host        = null;
    String              portStr     = null;
    int                 port        = 0;
    boolean             isIPv6      = false;
    Map<String,String>  props       = new TreeMap<String,String>(String.CASE_INSENSITIVE_ORDER);
    List<TGChannelUrl>  ftUrls      = null;


    private LinkUrl(String urlstr) throws TGException {

        this.urlstr = urlstr;
        parseInternal(urlstr.toLowerCase().trim());
    }

    private LinkUrl(Protocol protocol, String host, int port) {
        this.protocol = protocol;
        this.host = host;
        this.port = port;
        this.user = TGEnvironment.getInstance().getChannelDefaultUser();
        this.urlstr = null;

    }

    public static TGChannelUrl getDefaultUrl()  {
        return gDefaultUrl;
    }

    /**
     * URL is of the form protocol://host:value/{name=value;name=value...}
     *
     * @param url
     * @return
     * @throws Exception
     */
    public static TGChannelUrl parse(String url) throws  TGException {
        if ((url == null) || (url.length() == 0)) return gDefaultUrl;

        return new LinkUrl(url);
    }

    public static void main(String[] args) throws Exception {

        LinkUrl url = (LinkUrl) LinkUrl.parse("tcp://scott@foo.bar.com:8700");

        System.out.println("urlstr - case1 :" + url);

        url = (LinkUrl) LinkUrl.parse("tcp://foo.bar.com:8700/{userID=scott;ftHosts=foo1.bar.com,foo2.bar.com;sendSize=120}");

        System.out.println("urlstr - case2 :" + url);

        List<TGChannelUrl> ftUrls = url.getFTUrls();

        System.out.println("***** case 2: FT Urls....");
        for (TGChannelUrl tgurl : ftUrls) {
            System.out.println(tgurl);
            System.out.println("*****");
        }


        String[] cmd = {"/bin/sh", "-c", "stty raw </dev/tty"};
        Runtime.getRuntime().exec(cmd).waitFor();

        FileInputStream fin = new FileInputStream(FileDescriptor.in);
        System.setIn(fin);

        System.out.println("Enter the a character :");
        int b = System.in.read();
        char ch = (char) b;
        System.out.println("Char : " + ch);


        Console console = System.console();
        Reader reader = console.reader();
        ArrayList<Long> timeStamps = new ArrayList<Long>();
        StringBuilder password = new StringBuilder();
        timeStamps.add(System.currentTimeMillis());
        System.out.println("Enter your 8 character password");
        for (int i = 0; i < 8; i++) {
            password.append(reader.read());
            timeStamps.add(System.currentTimeMillis());
        }
        System.out.println("Timestamps: ");
        System.out.println(timeStamps);
        cmd = new String[]{"/bin/sh", "-c", "stty sane </dev/tty"};
        Runtime.getRuntime().exec(cmd).waitFor();

    }

    private void parseInternal(String url) throws TGException
    {
        String user =  parseProtocol(url);
        String hostAndPort = parseUser(user);
        String properties = parseHostAndPort(hostAndPort);
        parseProperties(properties.trim());
        //this.props.load(new StringReader(properties));
    }

    private String parseUser(String url) throws TGException
    {
        int posAt = url.indexOf('@');
        if (posAt == -1) return url;
        this.user = url.substring(0, posAt);
        props.put(ConfigName.ChannelUserID.getName(), user);
        return url.substring(posAt+1);
    }

    /**
     * Parse the Host and Port and return the remaing token.
     * @param hostAndPort
     * @return
     * @throws TGException
     */
    private String parseHostAndPort(String hostAndPort) throws TGException
    {
        if (hostAndPort.length() == 0) {
            this.host = "localhost";
            this.port = 8700;  //default values
            return null;
        }
        int offset = 0;
        int posIPv6 = hostAndPort.indexOf('[');
        if (posIPv6 != -1) {
            int endIPv6 = hostAndPort.indexOf(']');
            if (endIPv6 > posIPv6 + 2) {
                offset = endIPv6 + 1;
                this.isIPv6 = true;
            }
            else
                throw new TGException("Invalid or missing host name");
        }
        int pos = 0;
        if (this.isIPv6) {
            pos = hostAndPort.lastIndexOf(':');
        }
        else {
            pos = hostAndPort.indexOf(':', offset);
        }
        int lpos = hostAndPort.indexOf("/");
        if (pos < 0) {
            boolean noPort = true;
            if (hostAndPort.indexOf('.') < 0) {
                host = gDefaultHost;
                portStr = hostAndPort;
                noPort = false;
            }
            if (noPort) {
                host = hostAndPort;
                portStr = Integer.toString(gDefaultPort);
            }
        }
        else {
            int startpos = 0;
            int endpos = lpos != -1 ? lpos : hostAndPort.length();
            if (this.isIPv6) {
                startpos = 1;
                endpos = offset - 1;
            }
            host = hostAndPort.substring(startpos,pos);
            portStr = hostAndPort.substring(pos+1, endpos);
        }
        if (host.length() == 0) throw new TGException("Invalid or missing host name");
        if (portStr.length() == 0) throw new TGException("Invalid or missing port number");

        this.port = Integer.parseInt(portStr);
        return lpos == -1 ? "" : hostAndPort.substring(lpos+1);
    }

    /**
     * Parse the parse protocol, and return the remaining token
     * @param url
     * @return
     * @throws TGException
     */
    private String parseProtocol(String url) throws TGException
    {
        int pi;
        if (url.startsWith("tcp://")) {
            protocol = Protocol.TCP;
            pi = "tcp://".length();
        }
        else if (url.startsWith("ssl://")) {
            protocol = Protocol.SSL;
            pi = "ssl://".length();
        }
        else if (url.startsWith("http://")) {
            protocol = Protocol.HTTP;
            pi = "http://".length();
        }
        else if (url.startsWith("https://")) {
            protocol = Protocol.HTTPS;
            pi = "https://".length();
        }
        else {
            throw new TGException("Invalid protocol specification. for URL. URL is of the form protocol://host:value/{name=value;name=value...}");
        }
        return url.substring(pi);
    }

    /**
     * Parse the property string into a list of name value pairs. Note all names are in lower case.
     * @param properties
     * @return
     * @throws TGException
     */
    private String parseProperties(String properties) throws TGException
    {
        if (properties.length() == 0) return properties;

        if (!((properties.startsWith("{")) && (properties.endsWith("}")))) {
            throw new TGException("Malformed URL property specification - Must begin with { and end with }. All key=value must be seperated with ;");
        }

        String kvStr = properties.substring(1, properties.length()-1);
        char[] convtBuf = new char[1024];
        int limit;
        int keyLen;
        int valueStart;
        char c;
        boolean hasSep;
        boolean precedingBackslash;
        SemiColonReader scr = new SemiColonReader(kvStr.toCharArray());
        while ((limit = scr.readSemiColon()) > 0) {
            keyLen = 0;
            valueStart = limit;
            hasSep = false;

            //System.out.println("line=<" + new String(lineBuf, 0, limit) + ">");
            precedingBackslash = false;
            while (keyLen < limit) {
                c = scr.kvBuf[keyLen];
                //need check if escaped.
                if ((c == '=' || c == ':') && !precedingBackslash) {
                    valueStart = keyLen + 1;
                    hasSep = true;
                    break;
                } else if ((c == ' ' || c == '\t' || c == '\f') && !precedingBackslash) {
                    valueStart = keyLen + 1;
                    break;
                }
                if (c == '\\') {
                    precedingBackslash = !precedingBackslash;
                } else {
                    precedingBackslash = false;
                }
                keyLen++;
            }
            while (valueStart < limit) {
                c = scr.kvBuf[valueStart];
                if (c != ' ' && c != '\t' && c != '\f') {
                    if (!hasSep && (c == '=' || c == ':')) {
                        hasSep = true;
                    } else {
                        break;
                    }
                }
                valueStart++;
            }
            String key = parseUnicode(scr.kvBuf, 0, keyLen, convtBuf);
            String value = parseUnicode(scr.kvBuf, valueStart, limit - valueStart, convtBuf);
            ConfigName cn = ConfigName.fromName(key);
            if (cn == null) {
                props.put(key, value);
            }
            else {
                props.put(cn.getName(), value);
            }
        }
        return properties;
    }

    /*
     * Taken from loadProperties
     * Converts encoded &#92;uxxxx to unicode chars
     * and changes special saved chars to their original forms
     */
    private String parseUnicode(char[] in, int off, int len, char[] convtBuf)
    {
        if (convtBuf.length < len) {
            int newLen = len * 2;
            if (newLen < 0) {
                newLen = Integer.MAX_VALUE;
            }
            convtBuf = new char[newLen];
        }
        char aChar;
        char[] out = convtBuf;
        int outLen = 0;
        int end = off + len;
        while (off < end) {
            aChar = in[off++];
            if (aChar == '\\') {
                aChar = in[off++];
                if(aChar == 'u') {
                    // Read the xxxx
                    int value=0;
                    for (int i=0; i<4; i++) {
                        aChar = in[off++];
                        switch (aChar) {
                            case '0': case '1': case '2': case '3': case '4':
                            case '5': case '6': case '7': case '8': case '9':
                                value = (value << 4) + aChar - '0';
                                break;
                            case 'a': case 'b': case 'c':
                            case 'd': case 'e': case 'f':
                                value = (value << 4) + 10 + aChar - 'a';
                                break;
                            case 'A': case 'B': case 'C':
                            case 'D': case 'E': case 'F':
                                value = (value << 4) + 10 + aChar - 'A';
                                break;
                            default:
                                throw new IllegalArgumentException(
                                        "Malformed \\uxxxx encoding.");
                        }
                    }
                    out[outLen++] = (char)value;
                } else {
                    if (aChar == 't') aChar = '\t';
                    else if (aChar == 'r') aChar = '\r';
                    else if (aChar == 'n') aChar = '\n';
                    else if (aChar == 'f') aChar = '\f';
                    out[outLen++] = aChar;
                }
            } else {
                out[outLen++] = aChar;
            }
        }
        return new String (out, 0, outLen);
    }

    @Override
    public Protocol getProtocol() {
        return protocol;
    }

    @Override
    public String getHost()
    {
        if(host != null) return host;
        host = TGEnvironment.getInstance().getChannelDefaultHost();
        return host;
    }

    @Override
    public int getPort()
    {
        if (port != -1) return port;
        port = TGEnvironment.getInstance().getChannelDefaultPort();
        return port;
    }

    @Override
    public Map<String,String> getProperties() {
        return props;
    }

    public String getUser()
    {
        if (user != null) return user;
        user = props.get(ConfigName.ChannelUserID.getAlias());
        if ((user == null) || (user.length() == 0)) {
            user = props.get(ConfigName.ChannelUserID.getName());
            if ((user == null) || (user.length() == 0)) {
                user = TGEnvironment.getInstance().getChannelUser();
                if ((user == null) || (user.length() == 0)) {
                    user = TGEnvironment.getInstance().getChannelDefaultUser();
                }
            }
        }
        return user;
    }

    public synchronized List<TGChannelUrl> getFTUrls()
    {
        if (ftUrls != null) return ftUrls;

        String ftHosts = props.get(ConfigName.ChannelFTHosts.getAlias());
        if ((ftHosts == null) || (ftHosts.length() == 0)) {
            ftHosts = props.get(ConfigName.ChannelFTHosts.getName());
        }
        if ((ftHosts == null) || (ftHosts.length() == 0)) {
            ftUrls = new LinkedList<>();
            ftUrls.add(this);
            return ftUrls;
        }

        StringTokenizer st = new StringTokenizer(ftHosts.trim(), ",", false);
        ftUrls = new LinkedList<TGChannelUrl>();
        ftUrls.add(this);
        try {
            while (st.hasMoreTokens()) {
                String ftHost = st.nextToken();
                LinkUrl url = new LinkUrl(this.protocol, this.host, this.port);
                url.parseHostAndPort(ftHost);
                url.user = this.user;
                url.props = this.props;
                ftUrls.add(url);
            }
        }
        catch (Exception e) {
            e.printStackTrace(); //SS:TODO log
        }
        return ftUrls;
    }

    public String getUrlAsString() {
        if (this.urlstr != null) return this.urlstr;
        this.urlstr = String.format("%s://%s@%s:%d", protocol.name().toLowerCase(), user, host.toLowerCase(), port);
        return urlstr;
    }

    public String toString()
    {
        StringBuilder builder = new StringBuilder();
        builder.append("user:").append(getUser()).append(", ");
        builder.append("protocol:").append(getProtocol()).append(", ");
        builder.append("host:").append(getHost()).append(", ");
        builder.append("port:").append(getPort()).append(", ");
        builder.append("props:").append(getProperties());

        return builder.toString();
    }

    private class SemiColonReader {

        char inBuf[];
        int  curPos = 0;
        char kvBuf[];

        SemiColonReader(char[] inBuf) {
            this.inBuf = inBuf;
            kvBuf = new char[inBuf.length];
        }
        /**
         * Read upto a semicolon and return the nos of characters read.
         * @return
         */
        int readSemiColon()
        {
            int kvLen = 0;
            int kvPos = 0;
            boolean precedingBackslash = false;
            while (true)
            {
                if (curPos >= inBuf.length) return kvLen;
                char c = inBuf[curPos++];
                switch (c) {
                    case ' ':
                    case '\r':
                    case '\n':
                    case '\t':
                    case '\f':
                        break;
                    case '\\':
                        precedingBackslash = true;
                        break;
                    case ';':
                        if (precedingBackslash) {
                            kvBuf[kvPos++] = c;
                            ++kvLen;
                            precedingBackslash = false;
                            break;
                        }
                        else {
                            return kvLen;  //come out of the loop
                        }
                    default:
                        kvBuf[kvPos++] = c;
                        precedingBackslash = false;
                        ++kvLen;
                        break;
                }
            }
        }
    }
}
