#!/usr/bin/env bash


SERVER_NAME=fileserver


echo -e "os is $(uname -s)"

if [ "$(expr substr $(uname -s) 1 5)" == "MINGW" ];then
  SERVER_NAME=${SERVER_NAME}.exe
fi

rm -f ${SERVER_NAME}

go build -o ${SERVER_NAME}
