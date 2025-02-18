#!/bin/bash
set -e

while true
do
    peers=$(docker exec eirene0 bash -c "eirene attach /var/lib/eirene/data/eirene.ipc -exec 'admin.peers'")
    block=$(docker exec eirene0 bash -c "eirene attach /var/lib/eirene/data/eirene.ipc -exec 'eth.blockNumber'")

    if [[ -n "$peers" ]] && [[ -n "$block" ]]; then
        break
    fi
done

echo "$peers"
echo "$block"
