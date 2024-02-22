# ZetaChain Localnet Development & Testing Environment

ZetaChain localnet development and testing environment is divided into three different directories:

- [localnet](./contrib/localnet/README.md): a set of Docker images and script for spinning up a localnet.
- [e2e](./e2e/README.md): a set of Go packages for end-to-end testing between ZetaChain and other blockchains.
- [zetae2e](./cmd/zetae2e/README.md): a CLI tool using `e2e` for running end-to-end tests.

A description of each directory is provided in the respective README files.

## Running E2E tests

### Build zetanode
```bash
$ make zetanode
```

This Makefile rule builds the zetanode image. **Rebuild if zetacored/zetaclientd code is updated**.
```bash
$ docker build -t zetanode .
```

### Run e2e tests

Now we have built all the docker images, we can run the e2e test with make command:
```bash
make start-e2e-test
```

### Run upgrade tests

Upgrade tests run the E2E tests with an older version, upgrade the nodes to the new version, and run the E2E tests again.
This allows testing the upgrade process with a populated state.

Before running the upgrade tests, the versions must be specified in `Dockefile-upgrade`:

```dockerfile
ARG OLD_VERSION=v{old_major_version}.{old_minor_version}.{old_patch_version}
ENV NEW_VERSION=v{new_major_version}
```
The new version must match the version specified in `app/setup_handlers.go`

NOTE: We only specify the major version for `NEW_VERSION` since we use major version only for chain upgrade. Semver is needed for `OLD_VERSION` because we use this value to fetch the release tag from the GitHub repository.

The upgrade tests can be run with the following command:
```bash
make start-upgrade-test
```

The test the upgrade script faster a light version of the upgrade test can be run with the following command:
```bash
make start-upgrade-test-light
```
This command will run the upgrade test with a lower height and will not populate the state.

### Run stress tests

Stress tests run the E2E tests with a larger number of nodes and clients to test the performance of the network.
It also stresses the network by sending a large number of cross-chain transactions.

The stress tests can be run with the following command:
```bash
make start-stress-test
```

### Test logs

For all tests, the most straightforward logs to observe are the orchestrator logs.
If everything works fine, it should finish without panic.

The logs can be observed with the following command:
```bash
# in node/contrib/localnet/orchestrator
$ docker logs -f orchestrator
```

### Stop tests

To stop the tests,
```bash
make stop-test
```

### Run monitoring setup

Before starting the monitoring setup, make sure the Zetacore API is up at http://localhost:1317.
You can also add any additional ETH addresses to monitor in `zeta-node/contrib/localnet/grafana/addresses.txt` file

```bash
make start-monitoring
```

### Grafana credentials and dashboards

The Grafana default credentials are admin:admin. The dashboards are located at http://localhost:3000.

### Stop monitoring setup

```bash
make stop-monitoring
```

## Interacting with the Localnet

In addition to running automated tests, you can also interact with the localnet directly for more specific testing.

The localnet can be started without running tests with the following command:

```bash
make start-localnet
```

The localnet takes a few minutes to start. Printing the logs of the orchestrator will show when the localnet is ready. Once setup, it will display:
```
✅ the localnet has been setup
```

### Interaction with ZetaChain

ZetaChain
The user can connect to the `zetacore0` and directly use the node CLI with the zetacored binary with a funded account:

The account is named `operator` in the keyring and has the address: `zeta1amcsn7ja3608dj74xt93pcu5guffsyu2xfdcyp`

```bash
docker exec -it zetacore0 sh
```

Performing a query:

```bash
zetacored q bank balances zeta1amcsn7ja3608dj74xt93pcu5guffsyu2xfdcyp
```

Sending a transaction:

```bash
zetacored tx bank send operator zeta172uf5cwptuhllf6n4qsncd9v6xh59waxnu83kq 5000azeta --from operator --fees 2000000000000000azeta
```

### Interaction with EVM

The user can interact with the local Ethereum node with the exposed RPC on `http://0.0.0.0:8545`. The following testing account is funded:

```
Address: 0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC
Private key: d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263
```

Examples with the [cast](https://book.getfoundry.sh/cast/) CLI:

```bash
cast balance 0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC --rpc-url http://0.0.0.0:8545
98897999997945970464

cast send 0x9fd96203f7b22bCF72d9DCb40ff98302376cE09c --value 42 --rpc-url http://0.0.0.0:8545 --private-key "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263"
```

## Useful data

- TSS Address (on ETH): 0xF421292cb0d3c97b90EEEADfcD660B893592c6A2

## Add more e2e tests

The e2e tests are located in the e2e/e2etests package. New tests can be added. The process:

1. Add a new test file in the e2e/e2etests package, the `test_` prefix should be used for the file name.
2. Implement a method that satisfies the interface:
```go
type E2ETestFunc func(*E2ERunner)
```
3. Add the test to list in the `e2e/e2etests/e2etests.go` file.

The test can interact with the different networks using the runned object:
```go
type E2ERunner struct {
	ZevmClient   *ethclient.Client
	GoerliClient *ethclient.Client
	BtcRPCClient *rpcclient.Client

	CctxClient     crosschaintypes.QueryClient
	FungibleClient fungibletypes.QueryClient
	AuthClient     authtypes.QueryClient
	BankClient     banktypes.QueryClient
	ObserverClient observertypes.QueryClient
	ZetaTxServer   txserver.ZetaTxServer
	
	GoerliAuth *bind.TransactOpts
	ZevmAuth   *bind.TransactOpts
	
	// ...
}
```

## Localnet Governance Proposals

Localnet can be used for testing the creation and execution of governance propoosals.

Exec into the `zetacored0` docker container and run the script to automatically generate proposals in a variety of states and then extends the voting window to one hour, allowing you time to view a proposal in a pending state.
```
docker exec -it zetacore0 bash
/root/test-gov-proposals.sh
```
