#!/usr/bin/env bash


SERVER_NAME=fileserver

if [ "$(expr substr $(uname -s) 1 5)" == "MINGW" ]
then
  # Windows操作系统
  tskill ${SERVER_NAME}

  ./${SERVER_NAME}.exe >>error.log 2>&1 &

else
  # Linux、MACOS操作系统
  pkill ${SERVER_NAME}
  chmod +x ${SERVER_NAME}

  nohup ./${SERVER_NAME} >>error.log 2>&1 &

fi
