# Zeta Tool

Currently, has only one subcommand which finds inbound transactions or deposits that weren't observed on a particular
network. `filterdeposit`

## Configuring 

#### RPC endpoints
Configuring the tool for specific networks will require different reliable endpoints. For example, if you wanted to 
configure an ethereum rpc endpoint, then you will have to find an evm rpc endpoint for eth mainnet and set the field: 
`EthRPC`

#### Zeta URL
You will need to find an enpoint for zetachain and set the field: `ZetaURL`

#### TSS Addresses
Depending on which network you are using, you will have to populate the tss addresses for both EVM and BTC using these
fields: `TssAddressBTC`, `TssAddressEVM`

#### Contract Addresses
Depending on the network, connector and custody contract addresses must be set using these fields: `ConnectorAddress`,
`CustodyAddress`

#### EVM Block Ranges
When filtering evm transactions, a range of blocks is required and to reduce runtime of the command, a suitable range
must be selected and set in these fields: `EvmStartBlock`, `EvmMaxRange`

If a configuration file is not provided, a default config will be generated under the name 
`InboundTxFilter_config.json`. Below is an example of a configuration file used for mainnet: 

```json
{
 "ZetaURL": "http://46.4.15.110:1317",
 "TssAddressBTC": "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y",
 "TssAddressEVM": "0x70e967acfcc17c3941e87562161406d41676fd83",
 "BtcExplorer": "https://blockstream.info/api/address/bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y/txs",
 "EthRPC": "https://ethereum-rpc.publicnode.com",
 "ConnectorAddress": "0x000007Cf399229b2f5A4D043F20E90C9C98B7C6a",
 "CustodyAddress": "0x0000030Ec64DF25301d8414eE5a29588C4B0dE10",
 "EvmStartBlock": 19200110,
 "EvmMaxRange": 1000
}
```

## Running Tool

There are two targets available:

```
filter-missed-btc: install-zetatool
	./tool/filter_missed_deposits/filter_missed_btc.sh

filter-missed-eth: install-zetatool
	./tool/filter_missed_deposits/filter_missed_eth.sh
```

Running the commands can be simply done through the makefile in the node repo:

```
make filter-missed-btc
or ...
make filter-missed-eth
```