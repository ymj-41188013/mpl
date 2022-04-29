#!/bin/bash

function stop_gateway {
  # todo lab1-task-c
  # stop the gateway
  echo "implement me!"
  ID=`ps -ef | grep "sdbs_gateway" | grep -v "$0" | grep -v "grep" | awk '{print $2}'`
  echo $ID
  echo "---------------"
  for id in $ID
  do
  kill -9 $id
  echo "killed $id"
  done
  echo "---------------"
}

stop_gateway

