
# V21 Breaking Changes

### Update chain info refactored

* The `update_chain_info` message has been refactored to update/add a single chain ,instead of providing the entire list as `ChainInfo` 
  * The user is now required to provide a json file with the details of the chain to be updated/added.
    * If the chain already exists, the details will be updated.
    * If the chain does not exist, it will be added to the list of chains and saved to the store.
  * A new transaction type `RemoveChainInfo` has also been added to remove a chain from the list of chains.
    * It accepts the chain-id of the chain to be removed as a parameter.