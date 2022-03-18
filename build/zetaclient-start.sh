echo "Starting ZetaClient..."
echo $1 $2 $3 $4 $5

NODE_NUMBER=$1

NODE_0_ID=$2

NODE_0_IP=$3
 
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:/root/go/bin


## Start ZetaClient
FILE="/root/.zetaclient/chainobserver"
if  (( $NODE_NUMBER == 0 )) && [ -d "$FILE" ]; then
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
    zetaclientd  -peer /ip4/${NODE_0_IP}tcp/6668/p2p/${NODE_0_ID} -val val 2>&1 | tee ~/.zetaclient/zetaclient.log
elif (( $NODE_NUMBER > 0 )); then
    # Setup Zeta Client
    rm -f ~/.tssnew/address_book.seed 
    IDX=1 
    TSSPATH=/root/.tssnew 
    yes |  zetaclientd  -peer /ip4/${NODE_0_IP}/tcp/6668/p2p/${NODE_0_ID} -val val 2>&1 | tee ~/.zetaclient/zetaclient.log
fi



