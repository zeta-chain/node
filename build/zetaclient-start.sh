echo "Starting ZetaClient"
echo $1 $2 $3 $4 $5

NODE_NUMBER=$1
NODE_0_ID=$2
NODE_0_DNS=$3
 
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:/root/go/bin

echo "Waiting For Zetacored configuration"
sleep 5
## TODO -- Replace sleep. It should be able to determine if zetacored config has been completed or not

## Start ZetaClient
FILE="/root/.tssnew/e3234"
if  (( $NODE_NUMBER == 0 )) && [ -f "$FILE" ]; then
    echo "This is Node 0"
    echo "$FILE already exists."
    echo "Skipping ZetaClient Init"
    export MYIP=$(hostname -i)
    yes | zetaclientd -val val 2>&1 | tee ~/.zetaclient/zetaclient.log
elif (( $NODE_NUMBER == 0 )); then
    echo "This is Node 0"
    echo "Setting up ZetaClient For First Time Run"
    export MYIP=$(hostname -i)
    rm -f ~/.tssnew/address_book.seed 
    IDX=0 
    TSSPATH=/root/.tssnew 
    yes | zetaclientd -val val 2>&1 | tee ~/.zetaclient/zetaclient.log
fi


if  (( $NODE_NUMBER > 0 )) && [ -d "$FILE" ]; then
    echo "$FILE already exists."
    echo "Skipping ZetaClient Init"
    zetaclientd  -peer /dns/${NODE_0_DNS}tcp/6668/p2p/${NODE_0_ID} -val val 2>&1 | tee ~/.zetaclient/zetaclient.log
elif (( $NODE_NUMBER > 0 )); then
    # Setup Zeta Client
    rm -f ~/.tssnew/address_book.seed 
    IDX=1 
    TSSPATH=/root/.tssnew 
    yes |  zetaclientd  -peer /dns/${NODE_0_DNS}/tcp/6668/p2p/${NODE_0_ID} -val val 2>&1 | tee ~/.zetaclient/zetaclient.log
fi



