#!/bin/bash

function make_so {
  # todo
  # compile the plugin
  CGO_ENABLED=1 go build  -mod readonly --buildmode=plugin  -o codec.so ./codec.go
  echo "implement me!"
  exit -1
}

make_so

