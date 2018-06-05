#! /bin/bash

export dockerImage="tendermint/dbachain"
export dbaRoot="/root/.dba-test"

source ./libutils.sh

#######################
# Start Dba container
#######################
node_num=''
reinit=0

while getopts "hn:r" opt; do
    case $opt in
    h)
        printUsage
        exit 0
        ;;
    n)
        node_num=$OPTARG
        ;;
    r)
        echo "Important: Nodes will be reset before starting..."
        reinit=1
        rm -rf $dbaRoot
        ;;
    \?)
        echo "Invalid option: -$OPTARG"
        printUsage
        exit 0
        ;;
    esac
done

# Stop running containers
docker stop $(docker ps -a -q) > /dev/null 2>&1

# Build docker dba image
echo "Begin to build docker image..."
make build > /dev/null

echo "$node_num nodes will be started(include one seed)..."
seed_info=`startLocalMultiNodes $node_num $reinit`
