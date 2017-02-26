#!/bin/bash
set -e
rm sumoshell-linux.zip sumoshell-osx.zip

rm -rf /tmp/sumoshell-tmp
mkdir /tmp/sumoshell-tmp
cd /tmp/sumoshell-tmp
env GOOS=linux GOARCH=amd64 go build github.com/SumoLogic/sumoshell/sumo
env GOOS=linux GOARCH=amd64 go build github.com/SumoLogic/sumoshell/graph
env GOOS=linux GOARCH=amd64 go build github.com/SumoLogic/sumoshell/render
cd -
zip -rj sumoshell-linux.zip /tmp/sumoshell-tmp/*

set -e
rm -rf /tmp/sumoshell-tmp
mkdir /tmp/sumoshell-tmp
cd /tmp/sumoshell-tmp
env GOOS=darwin GOARCH=amd64 go build github.com/SumoLogic/sumoshell/sumo
env GOOS=darwin GOARCH=amd64 go build github.com/SumoLogic/sumoshell/graph
env GOOS=darwin GOARCH=amd64 go build github.com/SumoLogic/sumoshell/render
cd -
zip -rj sumoshell-osx.zip /tmp/sumoshell-tmp/*
