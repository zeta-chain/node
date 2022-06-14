#!/bin/bash
systemctl stop zetaclientd
systemctl stop zetacored
sleep 2 
echo "Stopped by CodeDeploy" >> /root/.zetacore/zetacored.log
echo "Stopped by CodeDeploy" >> /root/.zetaclient/zetaclient.log
