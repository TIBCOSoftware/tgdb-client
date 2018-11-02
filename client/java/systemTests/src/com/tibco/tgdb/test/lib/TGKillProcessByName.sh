#!/bin/sh
processName=$1
processKilled=false
OS=`uname`
export OS
case $OS in
 Darwin)
  processIds=`ps -ec | awk "/$processName/ "'{print $1}'`
  ;;
 *)
  processorType=`uname -p`
  export processorType
  case $processorType in
   s390x)
    processIds=`ps -ef | awk "/bin\/$processName/ "'{print $2}'`
	;;
   *)
    processIds=`ps -e | awk "/$processName/ "'{print $1}'`
	;;
  esac
  ;;
esac
for processId in $processIds
do
 kill -9 $processId
 processKilled=true
 echo Killed process $processName with pid $processId
done
if [ $processKilled = "false" ]; then
 echo "ERROR: The process \"$processName\" not found."
fi
