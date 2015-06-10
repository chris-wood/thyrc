#!/bin/bash

if [ -z "$GOPATH" ]; then
	export GOPATH=`pwd`
fi
go run client.go irc.freenode.net:6666

#package main

# go run client.go irc.freenode.net:6666

# startup protocol:
# PASS none
# NICK sorandom29      
# USER blah blah blah blah
