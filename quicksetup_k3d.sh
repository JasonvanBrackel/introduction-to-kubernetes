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

k3d cluster create hobbyfarm-example -p "80:80@loadbalancer"


# Without Zarf
kubectl create namespace hobbyfarm-system

helm install hobbyfarm hobbyfarm/hobbyfarm --namespace hobbyfarm-system --values ./values.yaml --wait --wait-for-jobs

echo "Waiting 8 seconds to wait for CRDs to settle"
sleep 8

kubectl apply -f manifests/admin-user.yaml -n hobbyfarm-system
kubectl apply -f manifests/hobbyfarm-admin-role.yaml
kubectl apply -f manifests/hobbyfarm-admin-rolebinding.yaml

# With Zarf
# export ARCH=arm64 # change this for your architecture

# zarf init

# zarf package create .

# zarf package deploy ./zarf-package-hobbyfarm-example-$ARCH.tar.zst