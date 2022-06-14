#!/bin/bash
systemctl stop zetaclientd
systemctl stop zetacored
echo "Stopped by CodeDeploy" >> /root/.zetacore/zetacored.log
echo "Stopped by CodeDeploy" >> /root/.zetaclient/zetaclient.log.log
sleep 2 
