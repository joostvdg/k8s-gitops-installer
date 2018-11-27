package core

var (
    Certificate = `apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  name: %s
  namespace: %s
spec:
  secretName: %s
  dnsNames:
  - %s
  acme:
    config:
    - http01:
        ingressClass: nginx
      domains:
      - %s
  issuerRef:
    name: %s
    kind: ClusterIssuer
`
)
// name, namespace, secretName, dns, dns, clusterIssueName
