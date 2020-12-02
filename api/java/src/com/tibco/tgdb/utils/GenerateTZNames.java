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

 * File name : GenerateTZNames.java
 * Created on: 09/22/2016
 * Created by: suresh
 * SVN Id: $Id: GenerateTZNames.java 3882 2020-04-17 00:27:45Z nimish $
 */

package com.tibco.tgdb.utils;

import java.util.TimeZone;

public class GenerateTZNames {
    static String FileHeader = "# \n" +
            "#  Copyright (c) 2016 TIBCO Software Inc.\n" +
            "#  All rights reserved.\n" +
            "#\n" +
            "#  File name: tznames.txt\n" +
            "#  Created on: September 15, 2016\n" +
            "#\n" +
            "#  SVN Id: $Id: GenerateTZNames.java 3882 2020-04-17 00:27:45Z nimish $\n" +
            "#\n" +
            "#  Lists all the timezone names with timezone ids. The timezone ids starts from 128.\n"+
            "#  Do not delete any lines or modify the timezone ids. \n" +
            "#  Timezones can be added, and they are done at the end.\n" +
            "#  The file format is as follows. \n" +
            "#  tzid tzname isactive isDST utcOffset aliases abbr \n"+
            "#  ------------------------------------------------  \n"+
            "#  tzid        - timezone id. \n" +
            "#  tzname      - timezone name listed as Area/Location\n" +
            "#  isactive    - [0|1]. Is the Timezone active. Occassionaly certain timezones are deactivated for political reasons\n"+
            "#  hasDST      - [0|1] Observes DST \n"+
            "#  utcOffset   - Utc offset specified as plus/minus from the UTC in hh:mm format\n" +
            "#  aliases     - A comma separate names with no spaces. Again political reasons, the location names are changes. \n" +
            "#                Example: Asia/Calcutta is changed to Asia/Kolkatta. Why - Political and historical reasons.\n" +
            "#                Sometime the classification for timezone is based on the most populated city. That can also change.\n" +
            "#  abbr        - The display name for readability. \n" +
            "#  ------------------------------------------------  \n";

    static String BeginHeader = "# **** BEGIN TIMEZONE GENERATION ******* !!! DO NOT EDIT BELOW !!! ";
    static String EndHeader = "# **** END TIMEZONE GENERATION ******* \n\n" +
                              "# **** BEGIN USER TIMEZONES **** !!! EDIT BELOW ONLY !!!\n" +
                              "# **** END USER TIMEZONES ****\n";


    public static void main(String[] args) {

        System.out.println(FileHeader);
        System.out.println(BeginHeader);
        String[] tznames = TimeZone.getAvailableIDs();
        int tzid = 128;
        for (String tzname:tznames) {
            TimeZone tz = TimeZone.getTimeZone(tzname);
            int rawoffset = tz.getRawOffset()/1000; //we need in secs.
            int hr = Math.abs(rawoffset / 3600);
            int mm = Math.abs((rawoffset % 3600)/60);
            System.out.printf("%d %-32s %d %d %c%02d:%02d\n",
                    tzid++,
                    tzname,
                    1,
                    tz.observesDaylightTime()?1:0,
                    rawoffset == 0 ? ' ':rawoffset > 0 ? '+' : '-',
                    hr,mm);

        }
        System.out.println(EndHeader);
        //ArrayBlockingQueue


    }





            

}
