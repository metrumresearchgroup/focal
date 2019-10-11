#!/bin/bash

#Change ownership of config files

chown -R focal /etc/focal/*

systemctl enable focal
systemctl daemon-reload
systemctl start focal