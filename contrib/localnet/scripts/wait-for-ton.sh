#!/bin/bash

timeout_seconds=300 # 5 minutes
poll_interval=2 # 2 seconds

ton_status_url="http://ton:8081/jsonRPC"

check_ton_status() {
    response=$(curl -s $ton_status_url -X POST \
        -H 'accept: application/json' -H 'Content-Type: application/json' \
        -d '{ "method": "getMasterchainInfo", "id": "1", "jsonrpc": "2.0" }')

    if [ -z "$response" ]; then
        echo "Waiting: no response"
        return 1
    fi

    if echo "$response" | jq -e '.ok == true' > /dev/null 2>&1; then
        echo "Pass: successful response received"
        return 0
    else
        echo "Waiting: unsuccessful response received"
        return 1
    fi
}

echo "💎 Checking TON status at $ton_status_url (timeout: $timeout_seconds seconds)"

elapsed=0
while [ $elapsed -lt $timeout_seconds ]; do
    if check_ton_status; then
        echo "💎 TON node bootstrapped"
        exit 0
    fi

    sleep $poll_interval
    elapsed=$((elapsed + poll_interval))
done

echo "💎 TON CHECK FAIL. Timeout reached ($timeout_seconds seconds)"
exit 1