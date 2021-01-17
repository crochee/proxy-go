@echo off
netstat -ano |findstr "8080"

tskill 6124   结束进程