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
 * File name : HttpChannel.${EXT}
 * Created on: 1/7/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: HttpChannel.java 583 2016-03-15 02:02:39Z vchung $
 */


package com.tibco.tgdb.channel.impl;


import com.tibco.tgdb.exception.TGAuthenticationException;
import com.tibco.tgdb.exception.TGException;

import java.io.*;
import java.util.Vector;

@Deprecated
/**
 *
 */
public class HttpChannel {


    private DataOutputStream    output;
    private DataInputStream     input;
    private LinkUrl             linkURL;
    private String              proxyUsername;
    private String              proxyPassword;

    private static class HttpStatus {
        String  HttpVersion = null;
        int     statusCode = 0;
        String  statusString = null;
    }

    private void doHttpConnect() throws IOException, TGException {
        // send connect string to proxy, telling it what SGDB server to connect to
        output.writeBytes("CONNECT "+linkURL.host+":"+linkURL.port+" HTTP/1.1\r\n");
        output.writeBytes("Host: "+linkURL.host+":"+linkURL.port+"\r\n");
        output.writeBytes("User-Agent: TIBCOSGDB/1.0\r\n");
        if (proxyUsername != null)
        {
            output.writeBytes("Proxy-Authorization: Basic " +
                    base64Encode(proxyUsername + ":" + proxyPassword) +
                    "\r\n");
        }
        output.writeBytes("\r\n");
        output.flush();

        Vector lines = readLines(input);
        HttpStatus status = processStatusLine((String) lines.firstElement());
        if (status.statusCode != 200)
        {
            if (status.statusCode == 407)
            {
                if (lines.size() < 2)
                    throw new TGException(
                            "Invalid HTTP response, not enough lines in response");

                // start from 1 because first element is status line
                for (int i=1; i < lines.size(); i++)
                {
                    String line = (String) lines.elementAt(i);
                    // if this finds the challenge it will throw the authentication exception,
                    // otherwise it just returns
                    processChallenge(line, status);
                }
                // couldn't parse the challenge
                throw new TGException("Invalid HTTP response, \"Basic\" challenge not found in response");
            }
            else
            {
                throw new TGException("Proxy rejected connection attempt " + status.statusString,  String.format("%d", status.statusCode));
            }
        }
    }

    static private Vector readLines(DataInputStream input) throws IOException, TGException {
        Vector lines = new Vector();

        BufferedReader reader = new BufferedReader(new InputStreamReader(input));
        String line = null;
        while ((line = reader.readLine()) != null && line.length() > 0)
        {
            // this is a continuation of the previous line, remove the old copy and add the new one
            if (line.charAt(0) == ' ' || line.charAt(0) == '\t')
            {
                if (lines.isEmpty())
                    throw new TGException("Invalid HTTP response, continuation line without previous line");
                String newLine = (String) lines.remove(lines.size()-1) + " " + line.trim();
                lines.add(newLine);
            }
            else
            {
                lines.add(line);
            }
        }

        return lines;
    }

    static private HttpStatus processStatusLine(String line) throws TGException {
        // as per the spec the status line must look like:
        //      HTTP/x.x <code>[ <reason>]
        // but we cope with proxies using more than one space between tokens
        HttpStatus status = new HttpStatus();

        line = line.trim();
        if (!line.startsWith("HTTP/"))
            throw new TGException(
                    "Invalid HTTP response, line does not start with \"HTTP/\": "+line);

        // get the version - HTTP/x.x
        int p = line.indexOf(' ');
        if (p <= 0)
            throw new TGException(
                    "Invalid HTTP response, unable to parse HTTP version: "+line);
        status.HttpVersion = line.substring(0, p).toUpperCase();

        // move past the spaces
        while (line.charAt(p) == ' ' || line.charAt(p) == '\t')
            p++;

        // get the status code - xxx
        int q = line.indexOf(' ', p);
        if (q < 0)
            q = line.length(); // Maybe no following text

        try {
            status.statusCode = Integer.parseInt(line.substring(p, q));
        }
        catch (NumberFormatException e) {
            throw new TGException(
                    "Invalid HTTP response, unable to parse status code: "+line);
        }

        // lastly get the reason, if there is one
        if (q < line.length())
            status.statusString = line.substring(q).trim();

        return status;
    }

    static private void processChallenge(String line, HttpStatus status)
            throws TGException {
        int p = line.indexOf(':');
        if (p > 0)
        {
            // only bother with the challenge header
            String header = line.substring(0, p).trim();
            if (header.equalsIgnoreCase("Proxy-Authenticate"))
            {
                // strip off the header name and ':', then trim the whitespace
                line = line.substring(p+1).trim();
                // this should throw the right exception
                processAuthenticate(line, status);
            }
        }
    }

    static private void processAuthenticate(String line, HttpStatus status)
            throws TGException
    {
        // there could be more than one challenge separated by whitespace and ','
        String realm = null;
        String scheme = null;
        do
        {
            String token = null;
            int p = line.indexOf(',');
            if (p > 0)
            {
                token = line.substring(0, p).trim();
            }
            else
            {
                token = line.trim();
                // set this so that else clause below returns ""
                p = line.length()-1;
            }

            // a  basic challenge is of the form Basic realm=xxxxxx
            int q = token.indexOf(' ');
            if (q > 0)
            {
                scheme = token.substring(0, q);
                if (scheme.equalsIgnoreCase("basic"))
                {
                    token = token.substring(q+1).trim();
                    q = token.indexOf('=');
                    realm = token.substring(q+1).trim();
                    if (realm.startsWith("\"") && realm.endsWith("\""))
                        realm = realm.substring(1, realm.length()-1);
                    throw new TGAuthenticationException(
                            "Proxy rejected connection attempt",
                            status.statusCode, status.statusString, realm);
                }
                else
                {
                    // ignore this token, move on
                    line = line.substring(p+1).trim();
                }
            }
            else
            {
                // could not find a spc character so this token is probably
                // just some parameter for a scheme we don't support,
                // move on
                line = line.substring(p+1).trim();
            }

        } while(line.length() > 0);

        throw new TGException("Proxy does not allow \"basic\" authentication");
    }

    static private final String codesymbols = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

    private String base64Encode(String s) {
        int len = s.length();
        int i;
        int rembytes = len%3;   // remainder
        // how many bytes we should add so total number is divisible by 3
        int addbytes = rembytes == 0 ? 0 : 3-rembytes;
        int charlen = ((len+addbytes)/3)*4;
        int triplets = len/3;  // number of triplets of bytes
        boolean lasttriplet = rembytes != 0; // if we have one extra triplet
        byte[] bytes = new byte[3];

        // +1 because first character specifies number of extra bytes
        char[] code = new char[charlen];

        // loop and encode every three characters as 4 characters
        int index;
        for (i=0; i<triplets; i++) {
            index = 4*i;
            bytes[0] = (byte)s.charAt(i*3);
            bytes[1] = (byte)s.charAt((i*3) + 1);
            bytes[2] = (byte)s.charAt((i*3) + 2);
            encode(code, index, bytes);
        }

        // encode last triplet if necessary
        if (lasttriplet) {
            for (i=0; i<rembytes; i++)
                bytes[i] = (byte)s.charAt(triplets*3+i);
            for (i=rembytes; i<3; i++)
                bytes[i] = 0;
            encode(code, triplets * 4, bytes);
        }

        String encoded = new String(code);

        return encoded;
    }

    private static void encode(char[] chars, int charoffset, byte[] bytes) {
        int k;
        k = (int) ((bytes[0]>>2)&0x3f);
        chars[charoffset] = (char)codesymbols.charAt(k);
        k = (int) (((bytes[0]<<4)&0x30) | ((bytes[1]>>4)&0x0f));
        chars[charoffset+1] = (char)codesymbols.charAt(k);
        if (bytes[1] != 0)
        {
            k = (int) (((bytes[1]<<2)&0x3c) | ((bytes[2]>>6)&0x03));
            chars[charoffset+2] = (char)codesymbols.charAt(k);
        }
        else
        {
            chars[charoffset+2] = '=';
        }
        if (bytes[2] != 0)
        {
            k = (int) (bytes[2]&0x3f);
            chars[charoffset+3] = (char)codesymbols.charAt(k);
        }
        else
        {
            chars[charoffset+3] = '=';
        }
    }

    // test program to test parsing of HTTP responses
    /*
    public static void main(String args[])
    {
        HttpStatus sts = null;
        int i;
        String[] statusLines = {"HTTP/1.1 200 OK",
                                "HTTP/1.0  \t  400 not ok ",
                                "HTTP/1.0    400 not ok ",
                                "HTTPS/2.3 200 OK",
                                " HTTP/1.1 200 OK", // this shouldn't happen because readLines will trim it
                                "HTTP/1.1 200",
                                "HTTP/1.1 407 Proxy Authentication Required"}; // end with a 407

        for (i = 0; i < statusLines.length; i++)
        {
            sts = null;
            try {
                System.out.println("Processing status line="+statusLines[i]);
                sts = processStatusLine(statusLines[i]);
            } catch (Exception e) {
                System.out.println("Exception="+e);
            }
            if (sts != null)
                System.out.println("Status version="+sts.HttpVersion+" code="+sts.statusCode+" string="+sts.statusString);
        }

        String[] input = {" \t invalid first line",
                          "HTTP/1.0 200 OK",
                          "Content-Type: text/html ",
                          "Content-Length:   \t   1354",
                          "Header1: some-long-value-1a, some-long-value-1b",
                          "Header1: \tsome-long-value-1a,",
                          "\t\t\t     some-long-value-1b",
                          "\t\t    \t some-long-value-2b",
                          "Proxy-Authenticate: Basic realm=\"foo realm\"",
                          "Proxy-Authenticate: Digest",
                          "                    realm=\"testrealm@host.com\",",
                          "                    qop=\"auth,auth-int\",",
                          "                    nonce=\"dcd98b7102dd2f0e8b11d0f600bfb0c093\",",
                          "                    opaque=\"5ccc069c403ebaf9f0171e9517f40e41\"",
                          "Proxy-Authenticate: Digest",
                          "                    realm=\"testrealm@host.com\",",
                          "                    qop=\"auth,auth-int\",",
                          "                    nonce=\"dcd98b7102dd2f0e8b11d0f600bfb0c093\",",
                          "                    opaque=\"5ccc069c403ebaf9f0171e9517f40e41\",",
                          "                    Basic realm=\"foo realm\"",
                          ""};

        String line = null;
        Vector lines = new Vector();
        i = 0;

        while (i < input.length && input[i].length() > 0)
        {
            line = input[i];
            // this is a continuation of the previous line, remove the old copy and add the new one
            if (line.charAt(0) == ' ' || line.charAt(0) == '\t')
            {
                if (lines.isEmpty()) {
                    System.out.println("Invalid HTTP response, continuation line without previous line");
                    i++;
                    continue;
                }
                String newLine = (String) lines.remove(lines.size()-1) + " " + line.trim();
                lines.add(newLine);
                i++;
            }
            else
            {
                lines.add(line);
                i++;
            }
        }
        for (i = 1; i < lines.size(); i++)
        {
            line = (String) lines.elementAt(i);
            System.out.println("Processing line: "+line);
            // if this finds the challenge it will throw the authentication exception,
            // otherwise it just returns
            try {
                processChallenge(line, sts);
            } catch (TGException e) {
                System.out.println("processChallenge threw exception: "+e);
            }
        }
    }
    */

}
