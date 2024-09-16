#!/bin/bash

timeout_seconds=300 # 5 minutes
poll_interval=10    # Check every 10 seconds

status_url="http://ton:8000/status"

check_status() {
    response=$(curl -s -w "\n%{http_code}" $status_url)
    body=$(echo "$response" | head -n 1)
    http_status=$(echo "$response" | tail -n 1)

    if [ "$http_status" == "200" ]; then
        echo "Pass: $body"
        return 0
    else
        echo "Waiting: $body"
        return 1
    fi
}

echo "💎 Checking TON status at $status_url (timeout: $timeout_seconds seconds)"

elapsed=0
while [ $elapsed -lt $timeout_seconds ]; do
    if check_status; then
        echo "💎 TON node bootstrapped"
        exit 0
    fi

    sleep $poll_interval
    elapsed=$((elapsed + poll_interval))
done

echo "💎 TON CHECK FAIL. Timeout reached ($timeout_seconds seconds)"
exit 1