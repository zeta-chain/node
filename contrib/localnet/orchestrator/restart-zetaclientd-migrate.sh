#!/bin/bash

echo restarting zetaclients

ssh -o "StrictHostKeyChecking no" "zetaclient0" -i ~/.ssh/localtest.pem killall zetaclientd
ssh -o "StrictHostKeyChecking no" "zetaclient1" -i ~/.ssh/localtest.pem killall zetaclientd
ssh -o "StrictHostKeyChecking no" "zetaclient0" -i ~/.ssh/localtest.pem "/usr/local/bin/zetaclientd start < /root/password.file > $HOME/zetaclient.log 2>&1 &"
ssh -o "StrictHostKeyChecking no" "zetaclient1" -i ~/.ssh/localtest.pem "/usr/local/bin/zetaclientd start < /root/password.file > $HOME/zetaclient.log 2>&1 &"

