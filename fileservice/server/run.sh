#!/usr/bin/env bash


SERVER_NAME=fileserver


echo -e "os is $(uname -s)"

if [ "$(expr substr $(uname -s) 1 5)" == "MINGW" ];then
  # Windows
  tskill ${SERVER_NAME}

  ./${SERVER_NAME}.exe >>error.log 2>&1 &

else
  # Linux、MACOS
  pkill ${SERVER_NAME}
  chmod +x ./${SERVER_NAME}

  nohup ./${SERVER_NAME} >>error.log 2>&1 &

fi
