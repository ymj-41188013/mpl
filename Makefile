SHELL = /bin/bash

lab0-task:
	make lab0-make-plugin
	cd ./test && GO111MODULE=on go test -v -run Task

lab0-task-a:
	cd ./test && GO111MODULE=on go test -v -run TaskA

lab0-task-b:
	cd ./test && GO111MODULE=on go test -v -run TaskB

lab0-task-c:
	make lab0-make-plugin
	cd ./test && GO111MODULE=on go test -v -run TaskC

lab0-make-plugin:
	cd ./pkg/plugin/demo && bash make_codec.sh
