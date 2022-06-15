#!/bin/bash
systemctl stop zetaclientd && echo "Stopped by CodeDeploy" >> /root/.zetaclient/zetaclient.log
systemctl stop zetacored && echo "Stopped by CodeDeploy" >> /root/.zetacore/zetacored.log
sleep 2 


