#!/bin/bash
set -e
GOOS=linux GOARCH=arm GOARM=5 go build main.go
echo "Enter password for ev3 ssh"
scp main robot@ev3dev.local:/home/robot
