GOOS=linux GOARCH=amd64 go build ./cmd/zetaclientd
echo "Copying to containers"
docker cp zetaclientd local-z-k8s-node0-1:/usr/local/bin/
docker cp zetaclientd local-z-k8s-node1-1:/usr/local/bin/
docker cp zetaclientd local-z-k8s-node2-1:/usr/local/bin/
docker cp zetaclientd local-z-k8s-node3-1:/usr/local/bin/
