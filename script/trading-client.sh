#!/bin/bash

# Mock 10 orders
for _ in $(seq 1 10)
do
    user_id=$(uuidgen)
    order_type=$(($RANDOM % 2))
    price_type=$(($RANDOM % 2))
    quantity=$(($RANDOM % 99 + 1))

    if [ $price_type -eq 0 ]
    then
        curl -X POST -d "{\"user_id\": \"$user_id\", \"order_type\": $order_type, \"quantity\": $quantity, \"price_type\": $price_type}" http://localhost:8080/order
    else
        price=$(($RANDOM % 10 + 95))
        curl -X POST -d "{\"user_id\": \"$user_id\", \"order_type\": $order_type, \"quantity\": $quantity, \"price_type\": $price_type, \"price\": $price}" http://localhost:8080/order
    fi
    echo
done
