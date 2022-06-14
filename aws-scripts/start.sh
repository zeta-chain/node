#!/bin/bash
echo "Started by CodeDeploy" >> /root/.zetacore/zetacored.log
echo "Started by CodeDeploy" >> /root/.zetaclient/zetaclient.log
systemctl start zetacored
systemctl start zetaclientd
