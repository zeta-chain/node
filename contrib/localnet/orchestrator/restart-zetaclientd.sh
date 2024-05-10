#!/bin/bash

# This script immediately restarts the zetaclientd on zetaclient0 and zetaclient1 containers in the network
# zetaclientd-supervisor will restart zetaclient automatically

echo restarting zetaclients

ssh -o "StrictHostKeyChecking no" "zetaclient0" -i ~/.ssh/localtest.pem killall zetaclientd
ssh -o "StrictHostKeyChecking no" "zetaclient1" -i ~/.ssh/localtest.pem killall zetaclientd