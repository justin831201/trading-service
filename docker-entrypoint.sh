#!/bin/sh

set -e

case "$@" in

run_trading_service)
    exec /srv/${PROJECT_NAME}/bin/${APP_NAME} -c /srv/${PROJECT_NAME}/config/config.yaml
    ;;

*)
    exec /bin/sh
    ;;

esac
