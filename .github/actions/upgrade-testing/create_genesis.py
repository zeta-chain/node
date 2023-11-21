import json
import os
from memory_profiler import profile


@profile
def genesis_modification():
    genesis = open(os.environ["NEW_GENESIS"], "r").read()
    genesis_json_object = json.loads(genesis)

    exported_genesis = open(os.environ["OLD_GENESIS"], "r").read()
    exported_genesis_json_object = json.loads(exported_genesis)

    exported_genesis = None
    genesis = None

    print("APP_STATE:AUTH:ACCOUNTS")
    genesis_json_object["app_state"]["auth"]["accounts"] = genesis_json_object["app_state"]["auth"]["accounts"] + \
                                                           exported_genesis_json_object["app_state"]["auth"]["accounts"]
    print("APP_STATE:CROSSCHAIN")
    genesis_json_object["app_state"]["crosschain"] = exported_genesis_json_object["app_state"]["crosschain"]
    print("APP_STATE:OBSERVER")
    genesis_json_object["app_state"]["observer"] = exported_genesis_json_object["app_state"]["observer"]
    print("APP_STATE:EMISSIONS")
    genesis_json_object["app_state"]["emissions"] = exported_genesis_json_object["app_state"]["emissions"]
    print("APP_STATE:FUNGIBLE")
    genesis_json_object["app_state"]["fungible"] = exported_genesis_json_object["app_state"]["fungible"]

    print("EVM:ACCOUNTS")
    evm_accounts = []
    for index, account in enumerate(exported_genesis_json_object["app_state"]["evm"]["accounts"]):
        if account["address"] == "0x0000000000000000000000000000000000000001":
            print("pop account", account["address"])
        elif account["address"] == "0x0000000000000000000000000000000000000006":
            print("pop account", account["address"])
        elif account["address"] == "0x0000000000000000000000000000000000000002":
            print("pop account", account["address"])
        elif account["address"] == "0x0000000000000000000000000000000000000002":
            print("pop account", account["address"])
        elif account["address"] == "0x0000000000000000000000000000000000000008":
            print("pop account", account["address"])
        else:
            evm_accounts.append(account)

    print("SET:EVM:ACCOUNTS")
    exported_genesis_json_object["app_state"]["evm"]["accounts"] = evm_accounts
    evm_accounts = None
    print("UPDATE:EVM")
    genesis_json_object["app_state"]["evm"] = exported_genesis_json_object["app_state"]["evm"]
    exported_genesis_json_object = None
    print("SAVE:EDITED:GENESIS")
    with open("genesis-edited.json", "w") as genesis:
        genesis.write(json.dumps(genesis_json_object))


genesis_modification()
