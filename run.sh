#!/bin/bash

if [ -z "$GOPATH" ]; then
	export GOPATH=`pwd`
fi
go run client.go irc.freenode.net:6666

