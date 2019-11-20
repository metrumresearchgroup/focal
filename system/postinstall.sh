#!/bin/bash

#Change ownership of focal configs

chown -R focal /etc/focal/*

#IS THIS SYSTEMD?
which systemctl

if [ $? -lt 1 ] ;
    then
    #this is a system D Box
    echo "This is a System D Service"
    #We've already copied the service file into the system. Let's reload, enable and start
    systemctl daemon-reload
    systemctl enable focal
    systemctl start focal
    exit 0 

    else
    #This is an upstart box
    echo "This is an init.d system"
    #Make sure the file is executable
    chmod +x /etc/init.d/focal
    update-rc.d focal defaults
    service focal start
    exit 0
fi