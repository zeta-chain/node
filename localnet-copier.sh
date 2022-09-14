rm -rf zetacored
rm -rf zetaclientd

echo "Building client"
GOOS=linux GOARCH=amd64 go build ./cmd/zetaclientd
echo "Copying to containers"
docker cp zetaclientd local-z-k8s-node0-1:/usr/local/bin/
docker cp zetaclientd local-z-k8s-node1-1:/usr/local/bin/
docker cp zetaclientd local-z-k8s-node2-1:/usr/local/bin/
docker cp zetaclientd local-z-k8s-node3-1:/usr/local/bin/

echo "Building core"
GOOS=linux GOARCH=amd64 go build ./cmd/zetacored
echo "Copying core to containers"
docker cp zetacored local-z-k8s-node0-1:/usr/local/bin/
docker cp zetacored local-z-k8s-node1-1:/usr/local/bin/
docker cp zetacored local-z-k8s-node2-1:/usr/local/bin/
docker cp zetacored local-z-k8s-node3-1:/usr/local/bin/
