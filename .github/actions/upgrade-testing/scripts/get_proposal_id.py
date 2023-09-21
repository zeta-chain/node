import json
import subprocess
import os

os.environ['NODE'] = "http://127.0.0.1:26657"
def run_command(self, cmd):
    COMMAND_PREFIX = "export PATH=" + self.go_path + ":${PATH} && "
    cmd = COMMAND_PREFIX + cmd
    result = subprocess.run(cmd, stdout=subprocess.PIPE, shell=True)
    result_output = result.stdout.decode('utf-8')
    return result_output

try:
    QUERY_GOV_PROPOSAL = f"""zetacored query gov proposals --output json --node {os.environ['NODE']}"""
    GOV_PROPOSALS = json.loads(run_command(QUERY_GOV_PROPOSAL))
    for proposal in GOV_PROPOSALS["proposals"]:
        try:
            PROPOSAL_ID = proposal["id"]
        except Exception as e:
            print(1)
    print(PROPOSAL_ID)
except Exception as e:
    print(1)