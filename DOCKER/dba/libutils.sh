#! /bin/bash

seedContainerID=''
seedNodeID=''
seedNodeIP=''

containerIDArray=()
nodeIDArray=()
nodeIPArray=()

persistent_peer=''
persistent_peers=''

reset=''

printUsage()
{
    echo "Dbachain Bootstrap Script

Usage:
    start_dba.sh -n num_of_nodes [-r]

Flags:
    -h  help for start_dba
    -n  number of nodes
    -r  unsafe reset all nodes before running
    "
}

startLocalMultiNodes()
{
    if (($# != 2));then
        echo "[startLocalMultiNodes] Input parameters is invalid! "
        return -1
    fi

    nodeNum=`expr $1 - 1`
    reset=`expr $2`

    for i in `seq 0 $nodeNum`;do
        node_info=`startLocalSingleNode $i`
        node_info_array=($node_info)
        containerID=${node_info_array[0]}
        nodeID=${node_info_array[1]}
        nodeIP=${node_info_array[2]}
        containerIDArray+=("$containerID")
        nodeIDArray+=("$nodeID")
        nodeIPArray+=("$nodeIP")

        if (($i == 0));then
            seedContainerID=$containerID
            seedNodeID=$nodeID
            seedNodeIP=$nodeIP
            persistent_peer="${seedNodeID}@${seedNodeIP}:46656"
            persistent_peers="$persistent_peer"
        else
            persistent_peer="${nodeID}@${nodeIP}:46656"
            persistent_peers="$persistent_peers,$persistent_peer"
        fi
    done

    return 0
}

startLocalSingleNode()
{
    if (($# != 1));then
        echo "[startLocalSingleNode] Input parameter is invalid! "
        return -1
    fi

    nodeNum=$1
    hostDir="/root/.dba-test/dba${nodeNum}"
    homeDir=$hostDir
    containerID=''
    nodeID=''
    nodeIP=''

    if (($reset == 1));then
        docker run -ti -v $hostDir:$homeDir $dockerImage /bin/bash -c \
            "dbachaind unsafe_reset_all --home $homeDir && dbachaind init --home $homeDir" > dba.log 2>&1

        # Modify node default config file setting
        sed -i 's/addr_book_strict = true/addr_book_strict = false/g' $homeDir/config/config.toml >> dba.log 2>&1
    fi

    if (($i == 0));then
        # Collect node info
        containerID=`docker run -d -v $hostDir:$homeDir $dockerImage /bin/bash -c \
            "dbachaind start --home $homeDir --p2p.seed_mode true"`
    else
        export seedChainID=`getChainID`

        cp $homeDir/config/genesis.json $homeDir/config/genesis.json.bak
        cat $homeDir/config/genesis.json.bak | jq 'to_entries |
            map(if .key == "chain_id"
                then . + {"value":env.seedChainID}
                else .
                end
            ) |
            from_entries' > $homeDir/config/genesis.json
        rm -f $homeDir/config/genesis.json.bak

        containerID=`docker run -d -v $hostDir:$homeDir $dockerImage /bin/bash \
            -c "dbachaind start --home $homeDir --p2p.persistent_peers ${persistent_peers}"`
    fi

    nodeID=`docker exec $containerID dbachaind show_node_id --home $homeDir`

    nodeIP=`docker inspect --format '{{ .NetworkSettings.IPAddress }}' $containerID`

    echo "$containerID $nodeID $nodeIP"

    return 0
}

getChainID()
{
    seedGenesisFile="$dbaRoot/dba0/config/genesis.json"
    if [ -e "$seedGenesisFile" ];then
        cat $seedGenesisFile | jq -r ".chain_id"
        return 0
    else
        cat $seedGenesisFile | jq -r ".chain_id"
        return -1
    fi
}


