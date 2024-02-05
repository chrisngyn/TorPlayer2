@echo off

go build -ldflags "-H windowsgui -extldflags=-static" -o build\bin\TorPlayer_%1.exe