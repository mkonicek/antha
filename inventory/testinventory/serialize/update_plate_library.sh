#! /bin/bash

go build
./serialize > ../platelibrary.go
go test
echo "Plate library updated"
