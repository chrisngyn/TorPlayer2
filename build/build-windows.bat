@echo off

go build -ldflags "-H windowsgui" -o build\bin\TorPlayer_%1.exe