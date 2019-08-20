#!/usr/bin/env bash


SERVER_NAME=fileserver


pkill ${SERVER_NAME}

chmod +x ${SERVER_NAME}

/usr/bin/nohup ./${SERVER_NAME} >> error.log 2>&1 &
