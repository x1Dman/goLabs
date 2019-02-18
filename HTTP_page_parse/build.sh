#!/bin/bash
export GOPATH=`pwd`
go install ./src/lenta
cd bin
./lenta
