#!/bin/bash
unset COLUMNS
PROG_NAME=$0
ACTION=$1
CURRENT_DIR=$(cd $(dirname $0);pwd)
SCRIPT_HOME=$(cd $(dirname $0)/..; pwd)
APP_NAME=cpts-build
SRC_RPM=${APP_NAME}-src
CONFIG_RPM=${APP_NAME}-config
SEC_RPM=${APP_NAME}-sec

deploy_user="service"
deploy_group="servicegroup"

usage() {
    echo "Usage: ${PROG_NAME} {start|stop|online|offline|install|uninstall|check_status}"
    exit 2
}

log() {
    echo "[$(date +"%Y-%m-%d %H:%M:%S")][${0##*/}:${FUNCNAME[1]}:${BASH_LINENO}] $*"
}

start()
{
    log "INFO: start ${APP_NAME} beginning."
    if [[ -f /opt/cloud/${APP_NAME}/bin/start.sh ]]; then
        su - $deploy_user -c "bash /opt/cloud/${APP_NAME}/bin/start.sh"
        [[ $? -ne 0 ]] && log "ERROR: start ${APP_NAME} failed." && exit 1
    else
        log "INFO: Doing noting about start."
    fi
    log "INFO: start ${APP_NAME} end."
}

stop()
{
    log "INFO: stop ${APP_NAME} beginning."
    if [[ -f /opt/cloud/${APP_NAME}/bin/stop.sh ]]; then
        su - $deploy_user -c "bash /opt/cloud/${APP_NAME}/bin/stop.sh"
        [[ $? -ne 0 ]] && log "ERROR: stop ${APP_NAME} failed." && exit 1
    else
        log "INFO: Doing noting about stop."
    fi
    log "INFO: stop ${APP_NAME} end."
}

online()
{
    log "INFO: online ${APP_NAME} beginning."
    # TODO Move it back when we have env-based appctl solution
    # sleep 5
    # systemctl start health_check
    log "INFO: online ${APP_NAME} end."
}

offline()
{
    log "INFO: offline ${APP_NAME} beginning."
    # TODO Move it back when we have env-based appctl solution
    # systemctl stop health_check
    # sleep 15
    log "INFO: offline ${APP_NAME} end."
}

install()
{
    log "INFO: install ${APP_NAME} beginning."
    log '==========rpm list ==========='
    if [[ -f ${CURRENT_DIR}/pre_install.sh ]]; then
        bash ${CURRENT_DIR}/pre_install.sh
        [[ $? -ne 0 ]] && log "ERROR: ${APP_NAME} pre_install failed." && exit 1
    fi

    rpm -ivh ${SCRIPT_HOME}/repo/${SRC_RPM}* --nodeps --force
    [[ $? -ne 0 ]] && log "ERROR: install ${APP_NAME} src rpm failed." && exit 2
    rpm -ivh ${SCRIPT_HOME}/config/${CONFIG_RPM}* --nodeps --force
    [[ $? -ne 0 ]] && log "ERROR: install ${APP_NAME} config rpm failed." && exit 3
    rpm -ivh ${SCRIPT_HOME}/bin/${BIN_RPM}* --nodeps --force
    [[ $? -ne 0 ]] && log "ERROR: install ${APP_NAME} bin rpm failed." && exit 4

    rpm -ivh ${SCRIPT_HOME}/sec/${SEC_RPM}* --nodeps --force
    [[ $? -ne 0 ]] && log "ERROR: install ${APP_NAME} sec rpm failed." && exit 4

    if [[ -f ${CURRENT_DIR}/post_install.sh ]]; then
        bash ${CURRENT_DIR}/post_install.sh
        [[ $? -ne 0 ]] && log "ERROR: ${APP_NAME} post_install failed." && exit 5
    fi

    chown $deploy_user:$deploy_group /opt/cloud/${APP_NAME}
    chmod 700 /opt/cloud/${APP_NAME}

    mkdir -p /opt/cloud/logs/${APP_NAME}
    chown $deploy_user:$deploy_group /opt/cloud/logs
    chmod 700 /opt/cloud/logs

    # start program on system boot
    # 实现cpts-build专用部署脚本
    if [[ -f /opt/cloud/${APP_NAME}/bin/install.sh ]]; then
        bash /opt/cloud/${APP_NAME}/bin/install.sh
        [[ $? -ne 0 ]] && log "ERROR: install ${APP_NAME} failed." && exit 1
    else
        log "INFO: Doing noting about install."
    fi

    log "INFO: install ${APP_NAME} end."
}

uninstall()
{
    # 实现cpts-build专用卸载脚本
    log "INFO: uninstall ${APP_NAME} beginning."
    if [[ -f /etc/init.d/${APP_NAME} ]]; then
        if [[ -e "/etc/redhat-release" ]]; then
            eulerOS=$(cat /etc/redhat-release | grep -i euler)
            if [[ -z "${eulerOS}" ]]; then
                chkconfig --del ${APP_NAME}
                rm -f /etc/init.d/${APP_NAME}
            fi
        fi
    fi

    if [[ -f ${CURRENT_DIR}/pre_uninstall.sh ]]; then
        bash ${CURRENT_DIR}/pre_uninstall.sh
        [[ $? -ne 0 ]] && log "ERROR: ${APP_NAME} pre_uninstall failed." && exit 1
    fi

    rpm -q ${SRC_RPM} | grep -v "not installed"| while read line
    do
        rpm -e ${line} > /dev/null 2>&1
    done

    rpm -q ${CONFIG_RPM} | grep -v "not installed"| while read line
    do
        rpm -e ${line} > /dev/null 2>&1
    done

    #Do not uninstall $SEC_RPM
    log "INFO: uninstall ${APP_NAME} end."
}

check_status()
{
    log "INFO: check ${APP_NAME} beginning."
    if [[ -f /opt/cloud/${APP_NAME}/bin/check_status.sh ]]; then
        su - $deploy_user -c "bash /opt/cloud/${APP_NAME}/bin/check_status.sh"
        [[ $? -ne 0 ]] && log "ERROR: check ${APP_NAME} status failed." && exit 1
    else
        log "INFO: Doing noting about check_status."
    fi
    log "INFO: check ${APP_NAME} end."
}

#check user
if [[ "${UID}" -ne 0 ]]; then
    log "ERROR: the script must run as root."
    exit 3
fi

#check usage
if [[ $# -lt 1 ]]; then
    usage
fi

case "${ACTION}" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    online)
        online
        ;;
    offline)
        offline
        ;;
    restart)
        offline
        stop
        start
        online
        ;;
    install)
        install
        ;;
    uninstall)
        uninstall
        ;;
    check_status)
        check_status
        ;;
    *)
        usage
        ;;
esac