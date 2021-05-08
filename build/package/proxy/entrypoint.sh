#!/bin/bash
set -ex
# docker build --no-cache -t obs .
# docker run -itd -p 8150:8150 -v /home/lcf/cloud/data:/obs/ --restart=always --name obsv1 obs
if [ "${1:0:1}" = '-' ]; then
  set -- proxy "$@"
fi

exec "$@"
