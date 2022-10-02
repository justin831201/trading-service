#!/bin/sh

set -e

case "$@" in

run_trading_service)
    echo "Wait for Infra Services..."
    sleep 10s
    echo "Start service..."
    exec /srv/${PROJECT_NAME}/bin/${APP_NAME} -c /srv/${PROJECT_NAME}/config/config.yaml
    ;;

*)
    exec /bin/sh
    ;;

esac
