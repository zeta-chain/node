#!/usr/bin/env bash

### chain init script for development purposes only ###

make clean install-zetacore
rm -rf ~/.zetacored
zetacored init test --chain-id=localnet -o

echo "Generating deterministic account - zeta"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | zetacored keys add zeta --recover --keyring-backend=test

echo "Generating deterministic account - mario"
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | zetacored keys add mario --recover --keyring-backend=test


zetacored add-genesis-account $(zetacored keys show zeta -a --keyring-backend=test) 500000000000000000000000000000000stake --keyring-backend=test
zetacored add-genesis-account $(zetacored keys show mario -a --keyring-backend=test) 50000000000000000000000000000000stake --keyring-backend=test

zetacored gentx zeta 1000000000000000000000000stake --chain-id=localnet --keyring-backend=test

echo "Collecting genesis txs..."
zetacored collect-gentxs

echo "Validating genesis file..."
zetacored validate-genesis
