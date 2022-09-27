#!/bin/sh

set -e

case "$@" in

run_trading_service)
    sleep 10s
    exec /srv/${PROJECT_NAME}/bin/${APP_NAME} -c /srv/${PROJECT_NAME}/config/config.yaml
    ;;

*)
    exec /bin/sh
    ;;

esac
