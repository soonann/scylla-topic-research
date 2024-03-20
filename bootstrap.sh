#!/usr/bin/env bash

# set the max processes
sudo sysctl -w fs.aio-max-nr=10485760
