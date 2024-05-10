#!/bin/bash

# Wait for authorized_keys file to exist (populated by zetacore0)
while [ ! -f ~/.ssh/authorized_keys ]; do
    echo "Waiting for authorized_keys file to exist..."
    sleep 1
done

while ! curl -s -o /dev/null zetacore0:26657/status ; do
    echo "Waiting for zetacore0 rpc"
    sleep 1
done

# copy zetacore0 keys
scp -R ~/.zetacored/config ~/.zetacored/os_info ~/.zetacored/config ~/.zetacored/keyring-file ~/.zetacored/

# serve binaries for upgrade proposal
python -m http.server