# ZetaClient Logging

- Log levels specified through zetaclientd init
    - `--log-level`   Flag
    - Log levels
        - `TRACE` (-1): for tracing the code execution path.
        - `DEBUG` (0): messages useful for troubleshooting the program.
        - `INFO` (1): messages describing the normal operation of an application.
        - `WARNING` (2): for logging events that need may need to be checked later.
        - `ERROR` (3): error messages for a specific operation.
        - `FATAL` (4): severe errors where the application cannot recover. `os.Exit(1)` is called after the message is logged.
        - `PANIC` (5): similar to `FATAL`, but `panic()` is called instead.

## Log Structure

- MasterLogger
    - StartupLogger : module = `Startup`
    - ZetaChainLogger
        - ChainLogger   : chain = `ZetaChain`
            - ZetaChainWatcher :  chain = `ZetaChain`   module=`ZetaChainWatcher`
    - BTCLogger
        - ChainLogger   : chain = `BTC`
            - WatchInbound  : chain = `BTC`   module=`WatchInbound`
            - WatchGasPrice : chain = `BTC`   module=`WatchGasPrice`
            - ObserverOutbound : chain = `BTC`  module=`ObserveOutbound`
            - WatchUTXOS:chain = `BTC`  module=`WatchUTXOS`
    - EVMLoggers ( Individual sections for each EVM Chain)
        - ChainLogger   : chain = `evm_chain_name`
            - BuildBlockIndex : chain = `evm_chain_name`   module=`BuildBlockIndex`
            - ExterrnalChainWatcher  : chain = `evm_chain_name`module=`ExternalChainWatcher`
            - WatchGasPrice : chain = `evm_chain_name`   module=`WatchGasPrice`
            - ObserverOutbound : chain = `evm_chain_name`  module=`ObserveOutbound`
    - BTCSigner : chain = `BTC`   module=`BTCsigner`
        - ProcessOutTX : chain = `BTC`   module=`BTCsigner`  OutboundId = `OuttxID of cctx being signed`  SendHash = `Index of cctx being signed`
    - EVMSigner : chain =  `evm_chain_name` module=`EVMSigner`
        - ProcessOutTX : chain =   `evm_chain_name` module=`BTCsigner`    OutboundId =  `OuttxID of cctx being signed` SendHash = `Index of cctx being signed`