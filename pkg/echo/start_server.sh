#!/system/bin/sh

BINARY="/data/local/bin/server"
LOGFILE="/data/local/tmp/server_script.log"

rm $LOGFILE
echo "Attempting to launch $BINARY" >> $LOGFILE
date >> $LOGFILE

$BINARY >> $LOGFILE 2>&1

EXIT_CODE=$?
echo "Binary launch command finished with exit code: $EXIT_CODE" >> $LOGFILE