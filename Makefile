SHELL = /bin/bash

lab0-task:
	make lab0-make-plugin
	cd ./test && GO111MODULE=on go test -v -run Lab0

lab0-task-a:
	cd ./test && GO111MODULE=on go test -v -run Lab0_TaskA

lab0-task-b:
	cd ./test && GO111MODULE=on go test -v -run Lab0_TaskB

lab0-task-c:
	make lab0-make-plugin
	cd ./test && GO111MODULE=on go test -v -run Lab0_TaskC

lab0-make-plugin:
	cd ./pkg/plugin/demo && bash make_codec.sh

lab1-task:
	cd ./test && GO111MODULE=on go test -v -run Lab1

lab1-task-a:
	cd ./test && GO111MODULE=on go test -v -run Lab1_TaskA

lab1-task-b:
	cd ./test && GO111MODULE=on go test -v -run Lab1_TaskB

lab1-task-c:
	cd ./test && GO111MODULE=on go test -v -run Lab1_TaskC
