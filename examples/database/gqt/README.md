TIBCO(R) Graph Database
Copyright (c) 2016-2021 TIBCO Software Inc. All rights reserved.
 
# GQT(Graph Query Traversal) Database

This file demonstrates how to import gqt data into TIBCO Graph 
database.

Instructions
------------------------------------------------------------------

Preliminary

Make sure you have the database/gqt directory copied to 
<tgdb_home>/examples/ 

Make sure to unzip import.zip file, after unzipping you should have 
<tgdb_home>/examples/gqt/import

<tgdb_home> specifies the path for the TGDB home directory. 
For Example,  
	On MacOSX /home/tibco/tgdb/3.1  
	On Windows C:/home/tibco/tgdb/3.1  

1. Adjust configuration files:
	* Copy tgdb-init.conf file into <tgdb_home>/bin directory
	
	* Open tgdb-init.conf and locate the [databases] section. 
	  Make sure that the gqtdb database is not commented out.
	  
	* Open ../examples/gqt/gqtdb.conf and locate the 
      [import] section. 
	  Make sure that this section is not commented out.
	  
	* Open <tgdb_home>/tgdb.conf and locate the databases section. 
	  Add gqtdb conf file path in this section.
	  For example,
	  gqtdb = ../examples/gqt/gqtdb.conf
	  
2. Init database and import GQT data:

	Navigate to directory <tgdb_home>/bin to execute any of the
	following commands. 

	Windows:
	tgdb -i -f -c tgdb-init.conf

	MacOSX/linux:
	./tgdb -i -f -c tgdb-init.conf

3. Start database

	Windows:
	tgdb -s -c tgdb.conf

	MacOSX/linux:
	./tgdb -s -c tgdb.conf


4. Launch an Admin console and connect to the server

	Navigate to <tgdb_home>/bin in a new command prompt window and execute following

	Windows:
	tgdb-admin

	MacOSX/linux:
	./tgdb-admin

	Execute a 'show types' command and make sure entries are present
	```
	admin@localhost:8223>show types
 	Name                  T SysId      #Entries
 	Default Nodetype      N 110        0
 	airlineType           N 9256       6162
 	airportType           N 9266       7184
 	allianceType          N 9258       6
 	cdi                   N 9270       421594
 	cdibatch              N 9264       63792
 	item                  N 9262       11115
 	lot                   N 9268       44774
 	machine               N 9254       148
 	workcenter            N 9260       104
 	Default Bidriecte...  E 1026       0
 	Default Directed ...  E 1025       0
 	Default Undirecte...  E 1024       0
 	contains              E 1041       736335
 	houses                E 1043       159
 	madefrom              E 1046       31666
 	makes                 E 1042       240726
 	memberType            E 1047       104
 	partof                E 1044       240421
 	routeType             E 1045       65600
 	stores                E 1049       421593
 	uses                  E 1048       981610
	22 types returned.
	```
