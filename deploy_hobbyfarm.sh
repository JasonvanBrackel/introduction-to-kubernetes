#!/bin/sh

set -e 

check_exists () {
    command -v "$1"
}

if [ $(check_exists "kubectl") ]   
then
    echo "kubectl found."
else
    echo "kubectl was not found.  You can find out more about installing kubectl at https://kubernetes.io/docs/tasks/tools/"
    exit 
fi

if [ $(check_exists "helm") ]   
then
    echo "helm found."
else
    echo "helm was not found.  You can find out more about installing helm at https://helm.sh/"
    exit 
fi

kubectl create namespace hobbyfarm-system

helm install hobbyfarm hobbyfarm/hobbyfarm --namespace hobbyfarm-system --values ./values.yaml --wait --wait-for-jobs

echo "Waiting 8 seconds to wait for CRDs to settle"
sleep 8

kubectl apply -f manifests/auth/admin-user.yaml -n hobbyfarm-system
kubectl apply -f manifests/auth/hobbyfarm-admin-role.yaml
kubectl apply -f manifests/auth/hobbyfarm-admin-rolebinding.yaml