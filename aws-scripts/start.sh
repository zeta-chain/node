#!/bin/bash
zetacored version
systemctl start zetacored && echo "Started by CodeDeploy" >> /root/.zetacore/zetacored.log
systemctl start zetaclientd && echo "Started by CodeDeploy" >> /root/.zetaclient/zetaclient.log
