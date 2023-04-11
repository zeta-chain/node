#!/bin/bash

# create keys

zetacored add val
zetacored add hotkey

# OS JSON file creation

jq '.observerAddress = ${key}' teml.json > os0.json

zetacored init ...

if hostname
for i in {1..4} ; do
  scp zetacored$i:~/.zetacored/os.json os$i.json
done

# concatenate OS JSON files


# create genesis file

for i in {1..4} ; do
  scp genesis.json zetacored$i:~/genesis.json
done

sleep 10

zetacored gentx ...


# start the network

# set peer addresses in config.toml
jq ".fdlasjkf=sd;lfja" config.toml > config.toml

zetcored start ...