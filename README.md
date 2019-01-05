# Kubernetes GitOps Installer

CLI for installing GitOps related resources in a Kubernetes cluster, much akin to [Jenkins X](jenkins-x.io).

The goal is not to compete with JenkinsX, but to learn how to build something like JenkinsX.

Currently, only GKE is supported.

## GKE Cluster

### Env variables

```bash
CLUSTER_NAME=joostvdg-2019-01-2
REGION=europe-west4
NODE_LOCATIONS=${REGION}-a,${REGION}-b
ZONE=europe-west4-a
K8S_VERSION=1.11.5-gke.4
```

### Check for supported cluster versions

```bash
gcloud container get-server-config --region $REGION
```

Only look at master node versions, those are usable in the `gcloud container clusters create` command.

### Create

```bash
gcloud container clusters create ${CLUSTER_NAME} \
    --region ${REGION} --node-locations ${NODE_LOCATIONS} \
    --cluster-version ${K8S_VERSION} \
    --num-nodes 2 --machine-type n1-standard-2 \
    --addons=HorizontalPodAutoscaling \
    --min-nodes 2 --max-nodes 3 \
    --enable-autoupgrade \
    --enable-autoscaling \
    --enable-network-policy \
    --labels=owner=jvandergriendt,purpose=practice
```

### Post Create

Before we can correctly install everything else, we need to have the cluster role `cluster-admin`.

```bash
kubectl create clusterrolebinding cluster-admin-binding \
    --clusterrole cluster-admin \
    --user $(gcloud config get-value account)
```

### Delete Cluster

```bash
gcloud container clusters delete $CLUSTER_NAME --region $REGION
```

### Configure Kubeconfig


#### Via GCloud

If you're using GCloud for authentication with GKE, you can use the following command to get your local kubeconfig configured.

```bash
gcloud container clusters get-credentials ${CLUSTER_NAME} --region ${REGION}
``` 

#### Via separate file

If you do not want to use GCloud, maybe for executing this in a container or CI server, you can create a new user/certificate.

This is taken from [gravitational](https://gravitational.com/blog/kubectl-gke/), which were nice enough to create a script for doing so.

You can find the script on [GitHub (get-kubeconfig.sh)](https://github.com/gravitational/teleport/blob/master/examples/gke-auth/get-kubeconfig.sh). 


## Install CloudBees Core

For installing CloudBees Core and other relevant tools.

```bash
go run main.go validate kubectl
# weavenet performance is pretty bad
# https://itnext.io/benchmark-results-of-kubernetes-network-plugins-cni-over-10gbit-s-network-36475925a560
# go run main.go install weavenet
sleep 10
go run main.go install helm
sleep 10
go run main.go install nginx
go run main.go gke sc-ssd
# Needs a pause
sleep 30
go run main.go install letsencrypt -e jvandergriendt@cloudbees.com --production
go run main.go cbc install --verbose --domainName cloudbees-core.kearos.net --production
go run main.go gke ing-ip -s ingress-nginx
```


## TODO

* allow weave-net to be encrypted or non-encrypted
* add vault integration
* add network policies
* add ldap
* add team recipes
* add a default team
