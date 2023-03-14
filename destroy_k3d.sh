#!/bin/sh

set -e 

check_exists () {
    command -v "$1"
}

if [ $(check_exists "k3d") ]   
then
    echo "k3d found."
else
    echo "k3d was not found.  You can find out more about installing k3d at k3d.io"
    exit 
fi

k3d cluster delete hobbyfarm-example