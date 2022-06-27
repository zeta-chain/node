#!/bin/bash
echo "echo $(zetacored version)" >> /root/.zetacore/zetacored.log
systemctl start zetacored && echo "Started by CodeDeploy" >> /root/.zetacore/zetacored.log
systemctl start zetaclientd && echo "Started by CodeDeploy" >> /root/.zetaclient/zetaclient.log
