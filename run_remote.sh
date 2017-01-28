#! /bin/bash
export GOPATH=$(pwd)
echo "Building executable"
go build main.go
echo "Start copying"
scp  -r ~/newHope/main madsjl@129.241.187.154:~/Desktop
echo "Finished copying"
ssh  madsjl@129.241.187.154
echo "Logged in"
