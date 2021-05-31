#!/bin/bash
set -ex
# docker build --no-cache -t obs .
# docker run -itd -p 8150:8150 -v /home/lcf/cloud/data:/obs/ --restart=always --name obsv1 obs
set -e

if [ "${1:0:1}" = '-' ]; then
  set -- proxy "$@"
fi

# ETH_0_IP=$(ifconfig eth0 | grep -w "inet" | awk '{print $2}')
# if [ -n "${ETH_0_IP}" ]; then
#   set -- "$@" --ip "${ETH_0_IP}"
# fi

if [[ "$1" = 'proxy' ]] ; then
  if [ "$(id -u)" = '0' ]; then
    echo "start with service"
    exec gosu serviceUser "$0" "$@"
  fi
fi
# 如果entrypoint.sh的入参在整个脚本中都没有被执行，那么exec "$@"会把入参执行一遍，
# 如果前面执行过了，这一行就不起作用
exec "$@"
