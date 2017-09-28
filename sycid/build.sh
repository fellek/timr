#!/bin/sh -e
env GOOS=linux GOARCH=arm GOARM=7 go build -o bin/sycid.read.arm  sycid.read.go
go build -o bin/sycid.read.x64  sycid.read.go