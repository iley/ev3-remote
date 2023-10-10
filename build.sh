#!/bin/bash
set -u -e -x
GOOS=linux GOARCH=arm GOARM=5 go build -mod=vendor
tar -czvf ev3-remote.tar.gz ev3-remote static
