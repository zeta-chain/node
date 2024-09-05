# TON localnet

This docker image represents a fully working TON node with 1 validator and a faucet wallet.

## How it works

- It uses [my-local-ton](https://github.com/neodix42/MyLocalTon) project without GUI.
  Port `4443` is used for [lite-client connection](https://docs.ton.org/participate/run-nodes/enable-liteserver-node).
- It also has a convenient sidecar on port `8000` with some useful tools.
- Please note that it might take **several minutes** to bootstrap the network.

## Sidecar

### Getting faucet wallet

```shell
curl -s http://ton:8000/faucet.json | jq
{
  "initialBalance": 1000001000000000,
  "privateKey": "...",
  "publicKey": "...",
  "walletRawAddress": "...",
  "mnemonic": "...",
  "walletVersion": "V3R2",
  "workChain": 0,
  "subWalletId": 42,
  "created": false
}
```

### Getting lite client config

Please note that the config returns IP of localhost (`int 2130706433`).
You need to replace it with the actual IP of the docker container.

```shell
curl -s http://ton:8000/lite-client.json | jq
{
  "@type": "config.global",
  "dht": { ... },
  "liteservers": [
    {
      "id": { "key": "...", "@type": "pub.ed25519" },
      "port": 4443,
      "ip": 2130706433
    }
  ],
  "validator": { ... }
}
```

### Checking the node status

It checks for config existence and the fact of faucet wallet deployment

```shell
curl -s http://ton:8000/status | jq
{
  "status": "OK"
}
```