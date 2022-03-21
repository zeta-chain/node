echo $0 $1

docker build -f ../Dockerfile.zetacore ../  -t zetacore
docker build -f ../Dockerfile.mockmpi ../ -t zeta-mockmpi
docker build -f ../Dockerfile.zetaclient ../  -t zetaclient