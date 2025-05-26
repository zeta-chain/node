## Usage 

### Show detailed statistics about the application database

### Command
```shell
zetatool db-stats --dbpath [pathToApplicationDb] --format [json|table]
```
### Example
```shell
zetatool db-stats --dbpath ./db --format json
{
  "totalKeys": 564492,
  "totalKeySize": 13540693,
  "totalValSize": 41652535,
  "modules": {
    "acc": {
      "count": 51796,
      "totalKeySize": 1087716,
      "totalValSize": 3868440,
      "avgKeySize": 21,
      "avgValSize": 74
    },
    "authority": {
      "count": 4363,
      "totalKeySize": 117801,
      "totalValSize": 59315,
      "avgKeySize": 27,
      "avgValSize": 13
    },
    "authz": {
      "count": 4389,
      "totalKeySize": 100947,
      "totalValSize": 61607,
      "avgKeySize": 23,
      "avgValSize": 14
    },
    "bank": {
      "count": 240816,
      "totalKeySize": 5297952,
      "totalValSize": 15499972,
      "avgKeySize": 22,
      "avgValSize": 64
    },
    "consensus": {
      "count": 4359,
      "totalKeySize": 117693,
      "totalValSize": 56754,
      "avgKeySize": 27,
      "avgValSize": 13
    },
    "crisis": {
      "count": 4359,
      "totalKeySize": 104616,
      "totalValSize": 56672,
      "avgKeySize": 24,
      "avgValSize": 13
    },
    "crosschain": {
      "count": 42507,
      "totalKeySize": 1190196,
      "totalValSize": 5299249,
      "avgKeySize": 28,
      "avgValSize": 124
    },
    "distribution": {
      "count": 57523,
      "totalKeySize": 1725690,
      "totalValSize": 3691463,
      "avgKeySize": 30,
      "avgValSize": 64
    },
    "emissions": {
      "count": 4584,
      "totalKeySize": 123768,
      "totalValSize": 94185,
      "avgKeySize": 27,
      "avgValSize": 20
    },
    "evidence": {
      "count": 4359,
      "totalKeySize": 113334,
      "totalValSize": 0,
      "avgKeySize": 26,
      "avgValSize": 0
    },
    "evm": {
      "count": 7395,
      "totalKeySize": 155295,
      "totalValSize": 575876,
      "avgKeySize": 21,
      "avgValSize": 77
    },
    "feemarket": {
      "count": 13077,
      "totalKeySize": 353079,
      "totalValSize": 610134,
      "avgKeySize": 27,
      "avgValSize": 46
    },
    "fungible": {
      "count": 4387,
      "totalKeySize": 114062,
      "totalValSize": 60447,
      "avgKeySize": 26,
      "avgValSize": 13
    },
    "gov": {
      "count": 4439,
      "totalKeySize": 93219,
      "totalValSize": 64217,
      "avgKeySize": 21,
      "avgValSize": 14
    },
    "group": {
      "count": 4363,
      "totalKeySize": 100349,
      "totalValSize": 56782,
      "avgKeySize": 23,
      "avgValSize": 13
    },
    "lightclient": {
      "count": 4359,
      "totalKeySize": 126411,
      "totalValSize": 56745,
      "avgKeySize": 29,
      "avgValSize": 13
    },
    "misc": {
      "count": 4360,
      "totalKeySize": 25055,
      "totalValSize": 4660205,
      "avgKeySize": 5,
      "avgValSize": 1068
    },
    "observer": {
      "count": 7939,
      "totalKeySize": 206414,
      "totalValSize": 488021,
      "avgKeySize": 26,
      "avgValSize": 61
    },
    "params": {
      "count": 4359,
      "totalKeySize": 104616,
      "totalValSize": 0,
      "avgKeySize": 24,
      "avgValSize": 0
    },
    "slashing": {
      "count": 13505,
      "totalKeySize": 351130,
      "totalValSize": 955079,
      "avgKeySize": 26,
      "avgValSize": 70
    },
    "staking": {
      "count": 72851,
      "totalKeySize": 1821275,
      "totalValSize": 5379169,
      "avgKeySize": 25,
      "avgValSize": 73
    },
    "upgrade": {
      "count": 4403,
      "totalKeySize": 110075,
      "totalValSize": 58203,
      "avgKeySize": 25,
      "avgValSize": 13
    }
  }
}
```

- `dbPath`: Path to application.db.
- `format`: Format of output, table (default) or json.
