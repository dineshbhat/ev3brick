#!/bin/bash
set -e
set +x
rm main || true
GOOS=linux GOARCH=arm GOARM=5 go build main.go
scp main robot@ev3dev.local:/home/robot
ssh robot@ev3dev.local './main'
