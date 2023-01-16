if [ $# -ne 1 ]
then 
  echo "Usage: node.sh <node num>"
  exit 1
fi
ssh -i z.pem -oStrictHostKeyChecking=no -p "702$1" alpine@0.0.0.0
