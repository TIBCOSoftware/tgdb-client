/**
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except
 *  in compliance with the License.
 *  A copy of the License is included in the distribution package with this file.
 *  You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 * Copyright (c) 2019 TIBCO Software Inc.
 * All rights reserved.
 *
 * File name : TGVersion.java
 * Created on: 2/15/19
 * Created by: suresh
 *
 * SVN Id: $Id $
 *
 *
 *  AUTOMATICALLY GENERATED AT BUILD TIME !!!!
 *
 *  DO NOT EDIT !!!
 */


package com.tibco.tgdb;




public class TGVersion {
    byte major;
    byte minor;
    byte update;
    byte hfno;
    short buildNo;
    byte buildType;
    byte edition;

    public static byte Production = 0;
    public static byte Engineering = 1;
    public static byte Beta = 2;

    public static byte Evaluation = 0;
    public static byte Community = 1;
    public static byte Enterprise = 2;
    public static byte Developer = 3;

    public static String PRODUCT_NAME  =   "TIBCO(R) Graph Database";
    public static String TIBCO_COPYRIGHT =  "Copyright (c) 2016-2019 TIBCO Software Inc. All rights reserved.";
    public static String TIBCO_LICENSE_STR1 =   "Please read the accompanying License and ReadMe documents;";
    public static String TIBCO_LICENSE_STR2 =   "your use of the software constitutes your acceptance of the terms contained in these documents.";

    static TGVersion gInstance = new TGVersion();

    TGVersion() {
        major       = @VERS_MAJOR@;
        minor       = @VERS_MINOR@;
        update      = @VERS_UPDATE@;
        hfno        = @VERS_HFNO@;
        buildNo     = @VERS_BUILDNO@;
        buildType   = @VERS_BUILDTYPE@;
        edition     = @VERS_EDITION@;
    }

	TGVersion(long version) 
	{
		major = (byte) (version & 0xff);
		minor = (byte) ((version & 0xff00) >> 8);
		update = (byte) ((version & 0xff0000) >> 16);
		hfno = (byte) ((version & 0xff000000) >> 24);
		buildNo = (short) ((version & 0xffff00000000L) >> 32);
		buildType = (byte) ((version & 0x0f000000000000L) >> 48);
		edition = (byte) ((version & 0xf0000000000000L) >> 52);
		
	}

    public static TGVersion getInstance()
    {
        return gInstance;
    }


    public long getAsLong() {
    	
    	long result = major;
    	long lMinor = (long)((long)minor << 8);
    	long lUpdate = (long)((long)update << 16);
    	long lhfNo = (long)((long)hfno << 24);
    	long lbuildNo = (long)((long)buildNo << 40);
    	long lbuildType = (long)((long)buildType << 44);
    	long lEdition = (long)((long)edition << 48);
    	
    	result = result | lMinor;
    	result = result | lUpdate;
    	result = result | lhfNo;
    	result = result | lbuildNo;
    	result = result | lbuildType; 
    	result = result | lEdition;
    	
        return result;
    }

    public String getLicense()
    {
        StringBuilder builder = new StringBuilder();
        //NH:TODO - See how the server side does.
        return builder.toString();
    }
    
    public byte getMajor() {
		return major;
	}

	public byte getMinor() {
		return minor;
	}

	public byte getUpdate() {
		return update;
	}
	
	public static TGVersion getInstanceFromLong(long version)
    {
        return new TGVersion(version);
    }
	
	@Override
	public boolean equals(Object obj) {

		if (obj instanceof TGVersion) 
		{
        	TGVersion versionObj = ((TGVersion)(obj));
        	if ((this.major == versionObj.major) && (this.minor == versionObj.minor) && (this.update == versionObj.update))
        	{
        		return true;
        	}
        }
        return false;
	}

	@Override
	public String toString() {
		return "TGVersion [major=" + major + ", minor=" + minor + ", update=" + update + "]";
	}
	
    
}