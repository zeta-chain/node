#!/bin/bash

java -jar my-local-ton.jar with-validators-1 nogui debug &

./sidecar &

# Wait for both processes to finish
wait -n
