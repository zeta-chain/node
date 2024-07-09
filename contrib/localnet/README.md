# ZetaChain Local Net Development & Testing Environment
This directory contains a set of scripts to help you quickly set up a 
local ZetaChain network for development and testing purposes. 
The scripts are based on the Docker 
and Docker Compose.

As a development testing environment, the setup aims to be
flexible, close to real world, and with fast turnaround
between edit code -> compile -> test results.

The `docker-compose.yml` file defines a network with:

* 2 zetacore nodes
* 2 zetaclient nodes
* 1 go-ethereum private net node (act as GOERLI testnet, with chainid 1337)
* 1 bitcoin core private node (planned; not yet done)
* 1 orchestrator node which coordinates E2E tests. 

The following Docker compose files can extend the default localnet setup:

- `docker-compose-stresstest.yml`: Spin up more nodes and clients for testing performance of the network.
- `docker-compose-upgrade.yml`: Spin up a network with with a upgrade proposal defined at a specific height.

Finally, `docker-compose-monitoring.yml` can be run separately to spin up a local grafana and prometheus setup to monitor the network.

## Running Localnet

Running the localnet requires `zetanode` Docker image. The following command should be run at the root of the repo:

```
make zetanode
```

Localnet can be started with Docker Compose:

```
docker-compose up -d
```

To stop the localnet:

```
docker-compose down
```

## Orchestrator

The `orchestrator` directory contains the orchestrator node which coordinates E2E tests. The orchestrator is responsible for:

- Initializing accounts on the local Ethereum network.
- Using `zetae2e` CLI to run the tests.
- Restarting ZetaClient during upgrade tests.

## Scripts

The `scripts` directory mainly contains the following scripts:

- `start-zetacored.sh`: Used by zetacore images to bootstrap genesis and start the nodes.
- `start-zetaclientd.sh`: Used by zetaclient images to setup TSS and start the clients.

## Prerequisites

The following are required to run the localnet:

- [Docker](https://docs.docker.com/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/downloads)
- [Go](https://golang.org/doc/install)
- [jq](https://stedolan.github.io/jq/download/)


### OP Integration

1 - Run localnet from `docker-compose-optimism.yml`.
>Adjust the exposed ports as necessary to avoid conflicts.

2 - Run localnet from OP documentation :  [Optimism Dev Node](https://docs.optimism.io/chain/testing/dev-node) 

3 - Make sure OP components are running alongside with Zeta localnet with existing EVM, BTC, and Zetachain nodes.

4 - Deploy OP smart contracts from [Optimism Smart Contract Deployment Guide](https://docs.optimism.io/builders/chain-operators/deploy/smart-contracts)

After running all scripts, you may encounter issues depending on your local environment, such as:

* Variables not being declared correctly. For instance, in my case, I had to modify `Config.sol` to deploy the contract. Some functions required static variables because this is just a local dev network.

>before:
```solidity
 function deployConfigPath() internal view returns (string memory _env) {
        if (vm.isContext(VmSafe.ForgeContext.TestGroup)) {
            _env = string.concat(vm.projectRoot(), "/deploy-config/hardhat.json");
        } else {
            _env = vm.envOr("DEPLOY_CONFIG_PATH", string(""));
            require(bytes(_env).length > 0, "Config: must set DEPLOY_CONFIG_PATH to filesystem path of deploy config");
        }
    }
 ```

 >after:
 ```solidity
 function deployConfigPath() internal pure returns (string memory _env) {
        _env = "deploy-config/devnetL1.json";
    }
 ```

* Error message like `(called 'Option::unwrap()' on a 'None' value)` which do indicates a problem inside the Foundry's forge tool, specifically within the `revm` library. 

To resolve this, I added debugging steps in my `deploy.sh` to ensure environment variables are set correctly.

````bash
echo "> Deploying contracts"
echo "RPC URL: $DEPLOY_ETH_RPC_URL"
echo "Private Key: $DEPLOY_PRIVATE_KEY"
echo "Config Path: $DEPLOY_CONFIG_PATH"
````

Ensure that Foundry and Forge are updated to the latest versions using  `foundryup`.

Compile a minimal Forge smart contract , such as `SimpleDeploy.s.sol`and execute it:

````solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleDeploy {
    function run() external {
        // Minimal deployment logic
    }
}
````

```sh
forge script -vvv scripts/deploy/SimpleDeploy.s.sol:SimpleDeploy --rpc-url 127.0.0.1:18545 --broadcast --private-key xxxxxxxxxxxxxxxxxxx
```

If execution is correct, try redepoying the set of contracts using `deploy.sh`.

- Deployment may proceed, but you might face issues with the `CREATE2` Deployer contract if it is not present in your local network: :  


```log
== Logs ==
  Writing artifact to /Users/aminechakrellah/Desktop/Projects/ZetaChain/optimism/packages/contracts-bedrock/deployments/1337-deploy.json
  Connected to network with chainid 1337
  Commit hash: 39bd919d8063ad57003eb2ac457743a241aa6df9
  DeployConfig: reading file deploy-config/devnetL1.json
  Deploying a fresh OP Stack including SuperchainConfig
  start of L1 Deploy!
  Deploying safe: SystemOwnerSafe with salt 0x47555c7c5eb40250af82c9713b290d445cad0893c01b18ae084f70d0b7b0d67d
  Saving SafeProxyFactory: 0x5FbDB2315678afecb367f032d93F642f64180aa3
  Saving SafeSingleton: 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
  Saving SystemOwnerSafe: 0xE97C3206DB3e054Ef065FF121BE168063721BA19
  New safe: SystemOwnerSafe deployed at 0xE97C3206DB3e054Ef065FF121BE168063721BA19
    Note that this safe is owned by the deployer key
  deployed Safe!
  Setting up Superchain
  Deploying AddressManager
  Saving AddressManager: 0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9
  AddressManager deployed at 0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9
  Deploying ProxyAdmin
  Saving ProxyAdmin: 0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9
  ProxyAdmin deployed at 0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9
  ProxyAdmin ownership transferred to Safe at: 0xE97C3206DB3e054Ef065FF121BE168063721BA19
  Deploying ERC1967 proxy for SuperchainConfigProxy
  Saving SuperchainConfigProxy: 0xa513E6E4b8f2a923D98304ec87F64353C4D5C853
     at 0xa513E6E4b8f2a923D98304ec87F64353C4D5C853
Error: 
script failed: missing CREATE2 deployer
```
Follow up with the [Optimism Documentation on how to deploy CREATE2 factory](https://docs.optimism.io/builders/chain-operators/tutorials/create-l2-rollup#deploy-the-create2-factory-optional).


- Add `rpc.allow-unprotected-txs` to geth if you cannot post transactions while deploying CREATE2.
- Use MetaMask to connect to your localnet to send transactions and fund smart contract addresses easily.
- Deploy the factory
```bash
cast publish --rpc-url http://0.0.0.0:18545 0xf8a58085174876e800830186a08080b853604580600e600039806000f350fe7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf31ba02222222222222222222222222222222222222222222222222222222222222222a02222222222222222222222222222222222222222222222222222222222222222 
```

>output
```json
{
  "status":"0x1",
  "cumulativeGasUsed":"0x10a23",
  "logs":[],
  "logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "type":"0x0",
  "transactionHash":"0xeddf9e61fb9d8f5111840daef55e5fde0041f5702856532cdbb5a02998033d26",
  "transactionIndex":"0x0",
  "blockHash":"0xa11e0bb5e8d2d6be633f16a46070fb1e2d8104778d84d1bb08fc997a3e754c26",
  "blockNumber":"0x2e2",
  "gasUsed":"0x10a23",
  "effectiveGasPrice":"0x174876e800",
  "from":"0x3fab184622dc19b6109349b94811493bf2a45362",
  "to":null,
  "contractAddress":"0x4e59b44847b379578588920ca78fbf26c0b4956c"
}
```

- Check Tx is mined : 

```json
{
  "jsonrpc":"2.0",
  "method":"eth_getTransactionReceipt",
  "params":["0xeddf9e61fb9d8f5111840daef55e5fde0041f5702856532cdbb5a02998033d26"],
  "id":1
}
```

- Verify the factory is deployed : 
```sh
cast codesize 0x4e59b44847b379578588920cA78FbF26c0B4956C --rpc-url http://0.0.0.0:18545
#output
69
```

At this stage you can redploy the L1 contracts :

```sh
npm run deploy
```

```log
==========================

ONCHAIN EXECUTION COMPLETE & SUCCESSFUL.

Transactions saved to: /Users/aminechakrellah/Desktop/Projects/ZetaChain/optimism/packages/contracts-bedrock/broadcast/Deploy.s.sol/1337/run-latest.json

Sensitive values saved to: /Users/aminechakrellah/Desktop/Projects/ZetaChain/optimism/packages/contracts-bedrock/cache/Deploy.s.sol/1337/run-latest.json
```

- Create genesis file : 

```sh
go run cmd/main.go genesis l2 \
  --deploy-config ../packages/contracts-bedrock/deploy-config/devnetL1.json \
  --l1-deployments ../packages/contracts-bedrock/deployments/devnetL1/.deploy \
  --outfile.l2 genesis.json \
  --outfile.rollup rollup.json \
  --l1-rpc 0.0.0.0:18545
```

















