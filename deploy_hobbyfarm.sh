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

helm repo add hobbyfarm https://hobbyfarm.github.io/hobbyfarm
helm repo add hf-shim-operator https://hobbyfarm.github.io/hf-shim-operator
helm repo add ec2-operator https://hobbyfarm.github.io/ec2-operator

# HobbyFarm an Infrastructure Automation
kubectl create namespace hobbyfarm
# AWS Credentials
helm install hobbyfarm hobbyfarm/hobbyfarm --namespace hobbyfarm --values ./values.yaml --wait --wait-for-jobs
helm install hf-shim-operator hf-shim-operator/hf-shim-operator --version v0.2.0-rc1 --namespace hobbyfarm --wait --wait-for-jobs
helm install ec2-operator ec2-operator/ec2-operator --version v0.2.0-rc1 --namespace hobbyfarm --wait --wait-for-jobs
# CRDS
kubectl apply -f https://raw.githubusercontent.com/hobbyfarm/metal-operator/main/chart/metal-operator/crds/equinix.cattle.io_importkeypairs.yaml
kubectl apply -f https://raw.githubusercontent.com/hobbyfarm/metal-operator/main/chart/metal-operator/crds/equinix.cattle.io_instances.yaml
kubectl apply -f https://raw.githubusercontent.com/hobbyfarm/metal-operator/main/chart/metal-operator/crds/importkeypair_crd.yaml
kubectl apply -f https://raw.githubusercontent.com/hobbyfarm/metal-operator/main/chart/metal-operator/crds/instance_crd.yaml
kubectl apply -f https://raw.githubusercontent.com/hobbyfarm/droplet-operator/master/charts/droplet-operator/crds/crd.yaml
# Auth
kubectl apply -f manifests/auth/admin-user.yaml -n hobbyfarm
kubectl apply -f manifests/auth/hobbyfarm-admin-role.yaml
kubectl apply -f manifests/auth/hobbyfarm-admin-rolebinding.yaml

kubectl create secret generic --from-literal=aws_access_key= --from-literal=aws_secret_key= aws-secret -n hobbyfarm

# An environment and sample course
kubectl apply -f courses/install-k3s/course.yaml
kubectl apply -f courses/install-k3s/environment.yaml
kubectl apply -f courses/install-k3s/scenario.yaml
kubectl apply -f courses/install-k3s/vmt.yaml
