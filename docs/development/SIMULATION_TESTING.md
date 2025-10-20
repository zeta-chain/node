# ZetaChain simulation testing
## Overview
The blockchain simulation tests how the blockchain application would behave under real life circumstances by generating 
and sending randomized messages.The goal of this is to detect and debug failures that could halt a live chain by 
providing logs and statistics about the operations run by the simulator as well as exporting the latest application
state.

## Simulation tests 

### Nondeterminism test
Nondeterminism test runs a full application simulation , and produces multiple blocks as per the config
It checks the determinism of the application by comparing the apphash at the end of each run to other runs
The test certifies that, for the same set of operations (regardless of what the operations are), we 
would reach the same final state if the initial state is the same
Approximate run time is 2 minutes.
```bash
make test-sim-nondeterminism
```

### Full application simulation test
Full application runs a full app simulation test with the provided configuration.
At the end of the run, it tries to export the genesis state to make sure the export works.
Approximate run time is 2 minutes.
```bash
make test-sim-full-app
```

### Import Export simulation test
The import export simulation test runs a full application simulation
and exports the application state at the end of the run.
This state is then imported into a new simulation.
At the end of the run, we compare the keys for the application state for both the simulations
to make sure they are the same.
Approximate run time is 2 minutes.
```bash
make test-sim-import-export
```

### Import and run simulation test
This simulation test exports the application state at the end of the run and imports it into a new simulation.
Approximate run time is 2 minutes.
```bash
make test-sim-after-import
```

### Multi seed long test
Multi seed long test runs a full application simulation with multiple seeds and multiple blocks.
It uses the `runsim` tool to run the same test in parallel threads.
Approximate run time is 30 minutes.
```bash
make test-sim-multi-seed-long
```

### Multi seed short test
Multi seed short test runs a full application simulation with multiple seeds and multiple blocks.
It uses the `runsim` tool to run the same test in parallel threads. 
This test is a shorter version of the Multi seed long test.
Approximate run time is 10 minutes.
```bash
make test-sim-multi-seed-short
```

### Import Export long test
This test runs the import export simulation test for a longer duration.
It uses the `runsim` tool to run the same test in parallel threads.
Approximate run time is 30 minutes.
```bash
make test-sim-import-export-long
```

### Import and run simulation test long
This test runs the import and run simulation test for a longer duration. 
It uses the `runsim` tool to run the same test in parallel threads.
Approximate run time is 30 minutes.
```bash
make test-sim-after-import-long
```

