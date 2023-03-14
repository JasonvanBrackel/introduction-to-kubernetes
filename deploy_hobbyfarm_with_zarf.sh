#!/bin/sh

set -e

if [ $(check_exists "zarf") ]   
then
    echo "zarf found."
else
    echo "zarf was not found.  You can find out more about installing zarf at https://zarf.dev"
    exit 
fi

export ARCH=arm64 # change this for your architecture

zarf init

zarf package create .

zarf package deploy ./zarf-package-hobbyfarm-example-$ARCH.tar.zst