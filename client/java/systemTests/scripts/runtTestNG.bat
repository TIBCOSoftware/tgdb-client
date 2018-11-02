@ECHO OFF
SET JAVA_HOME="C:\Program Files\Java\jdk1.8.0_74"
SET TGDB_HOME=C:\tgdb\1.0
SET TGDB_TEST=C:\svn\sgdb\trunk\qa\clientAPITest

%JAVA_HOME%\bin\java -cp %TGDB_TEST%/bin;%TGDB_TEST%/lib/testng-6.9.10/*;%TGDB_TEST%/lib/commons-exec-1.3/*;%TGDB_HOME%/lib/tgdb-client.jar -DTGDB_HOME=%TGDB_HOME%  org.testng.TestNG %*