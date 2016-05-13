package com.tibco.tgdb.utils;

import java.io.StringWriter;
import java.io.Writer;

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

 * File name :HexUtils
 * Created on: 12/26/14
 * Created by: suresh

 * SVN Id: $Id: HexUtils.java 622 2016-03-19 20:51:12Z ssubrama $
 */
public class HexUtils {
    
    static String NullString = "0000";
    static char Space = ' ';
    static String NewLine = "\r\n";
    
    public static String formatHex(byte[] buf)  {
        
        if (buf == null) return NullString;

        try (StringWriter sw = new StringWriter()) {

            formatHexToWriter(buf, sw, 0);
            return sw.toString();
        }
        catch (Exception e) {
            return NullString;
        }


        
    }
    
    public static String formatHex(byte[] buf, int actualLength)  {
        
        if (buf == null) return NullString;

        try (StringWriter sw = new StringWriter()) {

            formatHexToWriter(buf, sw, actualLength);
            return sw.toString();
        }
        catch (Exception e) {
            return NullString;
        }


        
    }

    public static void formatHexToWriter(byte[] buf, Writer writer, int actualLength) throws Exception
    {
        formatHexToWriter(buf, writer, 48, actualLength);
        return;
        
    }
    
    public static void formatHexToWriter(byte[] buf, Writer writer, int lineLength, int actualLength) throws Exception
    {

    	int blen = buf.length;
        boolean bNewLine = false;
        int lineNo = 1;
        

        writer.append("Formatted Byte Array:").append(NewLine);
        writer.append(String.format("%08x", 0));
        writer.append(Space);

        if (actualLength > 0) {
        	blen = actualLength;
        }
        for (int i=0; i < blen; i++)
        {
            if (bNewLine) {
                bNewLine = false;
                writer.append(NewLine);
                writer.append(String.format("%08x", (lineNo * lineLength) ));
                writer.append(Space);
            }

            writer.append(String.format("%02x", buf[i]));

            if ((i+1) % 2 ==0) writer.append(Space);

            if ((i+1) % lineLength == 0) {
                bNewLine = true;
                ++lineNo;
                writer.flush();
            }

        }
        
    }
}
