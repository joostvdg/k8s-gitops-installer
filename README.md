# k8s-gitops-installer

CLI for installing GitOps related resources in a Kubernetes cluster


```bash
CLUSTER_NAME=joostvdg-reg-nov18-3
REGION=europe-west4
NODE_LOCATIONS=${REGION}-a,${REGION}-b
ZONE=europe-west4-a
K8S_VERSION=1.11.2-gke.18
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
go run main.go install helm
go run main.go install nginx
go run main.go gke sc-ssd
go run main.go install letsencrypt -e jvandergriendt@cloudbees.com --production
go run main.go gke ing-ip
go run main.go cbc install --verbose --domainName cloudbees-core.kearos.net --production
```
