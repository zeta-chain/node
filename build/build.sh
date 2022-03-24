echo "You Entered $INPUT"

if  [ "$INPUT" == "" ]; then
    echo "Building zetacore, zetaclient, and zeta-mockpi"
    docker build -f ../Dockerfile.zetacore ../  -t zetacore
    docker build -f ../Dockerfile.mockmpi ../ -t zeta-mockmpi
    docker build -f ../Dockerfile.zetaclient ../  -t zetaclient
elif  [ "$INPUT" == "zetacore" ]; then
    echo "Building $INPUT Only"
    docker build -f ../Dockerfile.zetacore ../  -t zetacore
elif [ "$INPUT" == "zetaclient" ]; then
    echo "Building $INPUT Only"
    docker build -f ../Dockerfile.zetaclient ../  -t zetaclient
elif [ "$INPUT" == "zeta-mockmpi" ]; then
    echo "Building $INPUT Only"

    docker build -f ../Dockerfile.mockmpi ../ -t zeta-mockmpi
else 
    echo "Unknown Input"
    echo "Enter zetacore, zetaclient, zeta-mockmpi, or include no argument at all. If no argument is provided all three will be built."
fi
