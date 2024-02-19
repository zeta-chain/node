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
# in zeta-node/
$ docker build -t zetanode .
```

### Run e2e tests

Now we have built all the docker images; we can run the e2e test with make command:
```bash
# in zeta-node/
make start-e2e-test
```
which does the following docker compose command:
```bash
# in zeta-node/contrib/localnet/orchestrator
$ docker compose up -d
```

### Run upgrade tests

Upgrade tests run the E2E tests with an older version, upgrade the nodes to the new version, and run the E2E tests again.
This allows testing the upgrade process with a populated state.

Before running the upagrade tests, the versions must be specified in `Dockefile-upgrade`:

```dockerfile
ARG OLD_VERSION=vx.y.z
ENV NEW_VERSION=vxx.y.z
```
The new version must match the version specified in `app/setup_handlers.go`

The upgrade tests can be run with the following command:
```bash
# in zeta-node/
make start-upgrade-test
```

### Run stress tests

Stress tests run the E2E tests with a larger number of nodes and clients to test the performance of the network.
It also stresses the network by sending a large number of cross-chain transactions.

The stress tests can be run with the following command:
```bash
# in zeta-node/
make start-stress-test
```

### Test logs

For all tests, the most straightforward logs to observe are the orchestrator logs.
If everything works fine, it should finish without panic.

The logs can be observed with the following command:
```bash
# in zeta-node/contrib/localnet/orchestrator
$ docker logs -f orchestrator
```


### Stop tests

To stop the tests,
```bash
# in zeta-node/
make stop-test
```
which does the following docker compose command:
```bash
# in zeta-node/contrib/localnet/orchestrator
$ docker compose down --remove-orphans
```
### Run monitoring setup

Before starting the monitoring setup, make sure the Zetacore API is up at http://localhost:1317.
You can also add any additional ETH addresses to monitor in `zeta-node/contrib/localnet/grafana/addresses.txt` file

```bash
# in zeta-node/
make start-monitoring
```
which does the following docker compose command:
```bash
# in zeta-node/contrib/localnet/
$ docker compose -f docker-compose-monitoring.yml up -d
```
### Grafana credentials and dashboards

The Grafana default credentials are admin:admin. The dashboards are located at http://localhost:3000.

### Stop monitoring setup

```bash
# in zeta-node/
make stop-monitoring
```

which does the following docker compose command:

```bash
# in zeta-node/contrib/localnet/
$ docker compose -f docker-compose-monitoring.yml down --remove-orphans
```

## Useful data

- TSS Address (on ETH): 0xF421292cb0d3c97b90EEEADfcD660B893592c6A2

## Add more e2e tests

The e2e test (integration tests) are located in the
orchestrator/smoketest directory. The orchestrator is a Go program.

## Localnet Governance Proposals

Localnet can be used for testing the creation and execution of governance propoosals.

Exec into the `zetacored0` docker container and run the script to automatically generate proposals in a variety of states and then extends the voting window to one hour, allowing you time to view a proposal in a pending state.
```
docker exec -it zetacore0 bash
/root/test-gov-proposals.sh
```
