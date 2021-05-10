#!/bin/bash
set -ex
# docker build --no-cache -t obs .
# docker run -itd -p 8150:8150 -v /home/lcf/cloud/data:/obs/ --restart=always --name obsv1 obs
set -e

if [ "${1:0:1}" = '-' ]; then
  set -- proxy server "$@"
fi

#ETH_0_IP=$(ifconfig eth0 | grep -w "inet" | awk '{print $2}')
#if [ -n "${ETH_0_IP}" ]; then
#  set -- "$@" --ip "${ETH_0_IP}"
#fi

if [[ "$1" = 'proxy' ]] && [[ "$2" = 'server' ]]]; then
  if [ "$(id -u)" = '0' ]; then
    echo "start with service"
    exec gosu service "$0" "$@"
  fi
fi
exec "$@"
