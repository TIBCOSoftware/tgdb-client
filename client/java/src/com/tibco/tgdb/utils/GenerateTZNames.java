/**
 * Copyright (c) 2016 TIBCO Software Inc.
 * All rights reserved.
 * <p/>
 * File name : GenerateTZNames.${EXT}
 * Created on: 9/22/16
 * Created by: suresh
 * <p/>
 * SVN Id: $Id$
 */


package com.tibco.tgdb.utils;

import java.time.format.DateTimeFormatter;
import java.util.Date;
import java.util.TimeZone;
import java.util.concurrent.ArrayBlockingQueue;

public class GenerateTZNames {
    static String FileHeader = "# \n" +
            "#  Copyright (c) 2016 TIBCO Software Inc.\n" +
            "#  All rights reserved.\n" +
            "#\n" +
            "#  File name: tznames.txt\n" +
            "#  Created on: September 15, 2016\n" +
            "#\n" +
            "#  SVN Id: $Id:$\n" +
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
