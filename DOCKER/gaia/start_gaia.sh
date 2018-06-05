#! /bin/bash

# Currently gaia blockchain didn't support dynamic nodes depolyment, and it will be implemented
# by docker swarm mode later if needed.

export dockerImage="gaiatest"
export gaiaRoot="/root/.gaia-test"

export gen_txs_exist=true

source ./libutils.sh

while getopts "hr" opt; do
    case $opt in
    h)
        printUsage
        exit 0
        ;;
    r)
        echo "Important: Nodes will be reset before starting..."
        rm -rf $gaiaRoot
        gen_txs_exist=false
        ;;
    \?)
        echo "Invalid option: -$OPTARG"
        printUsage
        exit 0
        ;;
    esac
done

# Build gaia image
make build

# Stop running containers
docker stop $(docker ps -a -q) > /dev/null 2>&1

# Currently there are 4 nodes(include one seed) in the private gaia testnet, and it's for test only.
docker-compose down # To make sure that env is clean
docker-compose up




