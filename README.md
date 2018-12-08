# k8s-gitops-installer

CLI for installing GitOps related resources in a Kubernetes cluster

```bash
gcloud container get-server-config --region $REGION
```

```bash
CLUSTER_NAME=joostvdg-reg-dec18-1
REGION=europe-west4
NODE_LOCATIONS=${REGION}-a,${REGION}-b
ZONE=europe-west4-a
K8S_VERSION=1.11.3-gke.18
```

```bash
gcloud container clusters create ${CLUSTER_NAME} \
    --region ${REGION} --node-locations ${NODE_LOCATIONS} \
    --cluster-version ${K8S_VERSION} \
    --num-nodes 2 --machine-type n1-standard-2 \
    --addons=HorizontalPodAutoscaling \
    --min-nodes 2 --max-nodes 3 \
    --enable-autoupgrade \
    --enable-autoscaling \
    --labels=owner=jvandergriendt,purpose=practice

kubectl create clusterrolebinding cluster-admin-binding \
    --clusterrole cluster-admin \
    --user $(gcloud config get-value account)
go run main.go validate kubectl
go run main.go install weavenet
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


### Delete Cluster

```bash
gcloud container clusters delete $CLUSTER_NAME --region $REGION
```
