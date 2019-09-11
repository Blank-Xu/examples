#!/usr/bin/env bash


FILE_NAME=fileclient


echo -e "os is $(uname -s)"

if [ "$(expr substr $(uname -s) 1 5)" == "MINGW" ];then
  ${FILE_NAME}=${FILE_NAME}.exe
fi

rm -rf ${FILE_NAME}

go build -o ${FILE_NAME}
