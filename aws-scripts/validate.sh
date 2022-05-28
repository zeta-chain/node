#!/bin/bash

#zetacored
STATUS=$(systemctl is-active zetacored)
if [ "$STATUS" = "active" ]; then
    echo "Zetacored is running"
else
    echo "Zetacored is not running"
    exit 1
fi

#zetaclientd
STATUS=$(systemctl is-active zetaclientd)
if [ "$STATUS" = "active" ]; then
    echo "Zetaclientd is running"
else
    echo "Zetaclientd is not running"
    exit 1
fi
