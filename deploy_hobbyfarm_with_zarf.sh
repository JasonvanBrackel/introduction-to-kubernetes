#!/bin/sh

set -e

check_exists () {
    command -v "$1"
}

if [ $(check_exists "zarf") ]   
then
    echo "zarf found."
else
    echo "zarf was not found.  You can find out more about installing zarf at https://zarf.dev"
    exit 
fi

export ARCH=arm64 # change this for your architecture

zarf init --confirm

zarf package create . --confirm

zarf package deploy ./zarf-package-hobbyfarm-example-$ARCH.tar.zst --confirm