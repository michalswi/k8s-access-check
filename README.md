# k8s-access-check

Related to checking API access described [here](https://kubernetes.io/docs/reference/access-authn-authz/authorization/#checking-api-access).

## Run
```sh
# parameters (check './main -h')
--dir <path_to_json> 
--ns <namespace> 
--run-outside-k-cluster true

# admin example
$ ./main --dir template.json --ns kube-system --run-outside-k-cluster true
2019/08/02 16:07:44 Init namespace: kube-system
2019/08/02 16:07:44 K8s version: v1.13.5
2019/08/02 16:07:44 User can /create/ a /deployments/, status: ALLOWED
2019/08/02 16:07:44 User can /create/ a /pods/, status: ALLOWED
2019/08/02 16:07:44 User can /get/ a /deployments/, status: ALLOWED
2019/08/02 16:07:44 User can /get/ a /pods/, status: ALLOWED
```

## Test Service Account
```sh
$ kubectl apply -f rbac.yml

$ kubectl --as=system:serviceaccount:default:michal -n kube-system get deployments

# test SA permission in the default namespace
$ ./main --dir template.json --run-outside-k-cluster true
2019/08/02 16:09:53 Init namespace: default
2019/08/02 16:09:53 K8s version: v1.13.5
2019/08/02 16:09:53 User can /create/ a /deployments/, status: ALLOWED
2019/08/02 16:09:53 User can /create/ a /pods/, status: NOTALLOWED
2019/08/02 16:09:54 User can /get/ a /deployments/, status: ALLOWED
2019/08/02 16:09:54 User can /get/ a /pods/, status: NOTALLOWED

```