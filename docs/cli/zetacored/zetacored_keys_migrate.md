# keys migrate

Migrate keys from amino to proto serialization format

### Synopsis

Migrate keys from Amino to Protocol Buffers records.
For each key material entry, the command will check if the key can be deserialized using proto.
If this is the case, the key is already migrated. Therefore, we skip it and continue with a next one. 
Otherwise, we try to deserialize it using Amino into LegacyInfo. If this attempt is successful, we serialize 
LegacyInfo to Protobuf serialization format and overwrite the keyring entry. If any error occurred, it will be 
outputted in CLI and migration will be continued until all keys in the keyring DB are exhausted.
See https://github.com/cosmos/cosmos-sdk/pull/9695 for more details.

It is recommended to run in 'dry-run' mode first to verify all key migration material.


```
zetacored keys migrate [flags]
```

### Options

```
  -h, --help   help for migrate
```

### Options inherited from parent commands

```
      --home string              The application home directory 
      --keyring-backend string   Select keyring's backend (os|file|test) 
      --keyring-dir string       The client Keyring directory; if omitted, the default 'home' directory will be used
      --log_format string        The logging format (json|plain) 
      --log_level string         The logging level (trace|debug|info|warn|error|fatal|panic) 
      --output string            Output format (text|json) 
      --trace                    print out full stack trace on errors
```

### SEE ALSO

* [zetacored keys](zetacored_keys.md)	 - Manage your application's keys

