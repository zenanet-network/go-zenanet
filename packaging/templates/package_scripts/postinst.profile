#!/bin/bash
# This is a postinstallation script so the service can be configured and started when requested
#
if [ -d "/var/lib/eirene" ]
then
    echo "Directory /var/lib/eirene exists."
else
    mkdir -p /var/lib/eirene
    sudo chown -R eirene /var/lib/eirene
fi
sudo systemctl daemon-reload
