#!bin/bash
set -ex
export enable_log=true
export log_path=./log/proxy.log
export log_level=DEBUG
export GIN_MODE=release
export config=./conf/config.yml
#./proxy >/dev/null &
./proxy >>./log/console.txt 2>&1 &
