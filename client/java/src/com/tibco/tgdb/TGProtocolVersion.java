/**
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
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
 * <p>
 * File name :TGProtocolVersion
 * Created on: 12/24/14
 * Created by: suresh
 * </p>
 * SVN Id: $Id: TGProtocolVersion.java 3123 2019-04-25 21:48:54Z nimish $
 */

package com.tibco.tgdb;

public class TGProtocolVersion {

    final static byte MAJOR_VERSION = 1;
    final static byte MINOR_VERSION = 0;
    private final static int TGMAGIC = 0xdb2d1e4;

    public static short  getProtocolVersion() {
        short ver = (int) MAJOR_VERSION;
        ver = (short)(ver << 8);
        ver = (short)(ver + MINOR_VERSION);

        return ver;
    }

    public static int getMagic() {
        return TGMAGIC;
    }

    //SS:TODO : Version Compatibility Checks needs to be done.
    //Client protocol could be higher than the Server.
    public static boolean isCompatible(int protocolVersion)
    {
        return protocolVersion == getProtocolVersion();
    }
}
