# Zetachain simulation testing
## Overview
The blockchain simulation tests how the blockchain application would behave under real life circumstances by generating 
and sending randomized messages.The goal of this is to detect and debug failures that could halt a live chain,by providing 
logs and statistics about the operations run by the simulator as well as exporting the latest application state.


## Simulation tests 

### Nondeterminism test
Nondeterminism test runs a full application simulation , and produces multiple blocks as per the config
It checks the determinism of the application by comparing the apphash at the end of each run to other runs
The test certifies that , for the same set of operations ( irrespective of what the operations are ), we would reach the same final state if the initial state is the same
```bash
make test-sim-nondeterminism
```
### Full application simulation test
Full application runs a full app simulation test with the provided configuration.
At the end of the run it tries to export the genesis state to make sure the export works.
```bash
make test-sim-full-app
```

### Multi seed long test
Multi seed long test runs a full application simulation with multiple seeds and multiple blocks.This runs the test for a longer duration compared to the multi seed short test
```bash
make test-sim-multi-seed-long
```

### Multi seed short test
Multi seed short test runs a full application simulation with multiple seeds and multiple blocks. This runs the test for a longer duration compared to the multi seed long test
```bash
make test-sim-multi-seed-short
```