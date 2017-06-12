#### Maya API server and Kubernetes

Maya API server is launched as a deploy/pod unit in Kubernetes. This service is 
the interface for the storage clients to operate on OpenEBS storage. Typically, 
volume plugins (here K8s Flex Volume driver) act as http clients to Maya API 
service. 

> OpenEBS has the concept of VSM (Volume Storage Machine) to provide persistent
storage. Maya API service provides operations w.r.t VSM as a unit.

##### Launch Maya API server as a K8s Pod

- `Create the yaml specs for launching Maya API server as a K8s Pod`

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: maya-apiserver
spec:
  containers:
  - image: openebs/m-apiserver:test
    imagePullPolicy: IfNotPresent
    name: maya-apiserver
    ports:
    - containerPort: 5656
```

- `Use kubectl to launch Maya API server as a K8s Pod`

```bash
kubectl create -f maya-api-server.yaml
```

- `Get the IP address of above created Pod`

```bash
kubectl describe pod/maya-apiserver
```

- `Create yaml specs for role associated with VSM operations`

```
$ cat vsm-role-all.yaml
# This role allows operations on K8s pods, deployments in "default" namespace
kind: Role
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: vsm-role-all
  namespace: default
rules:
- apiGroups: ["apps", "extensions"]
  resources: ["pods","deployments"]
  verbs: ["get","list","watch","create","update","patch","delete"]
```

```
$ cat vsm-service.yaml
# This role binding binds vsm-role-all role to default service account
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata: 
  name: vsm-service
  namespace: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: vsm-role-all
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: system:serviceaccount:default:default
```

```
$ cat vsm-cluster-role.yaml
# This cluster role allows operations on K8s services in "default" namespace
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  namespace: default
  name: vsm-cluster-role
rules:
- apiGroups: ["*"]
  resources: ["services"]
  verbs: ["*"]
```

- `Create yaml specs for cluster role binding w.r.t VSM operations`

```
$ cat vsm-cluster-service.yaml
# This role binding allows "default" service account to bind to 
# vsm-cluster-role in "default" namespace.
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: vsm-cluster-service
  namespace: default
subjects:
- kind: User
  name: system:serviceaccount:default:default
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: vsm-cluster-role
  apiGroup: rbac.authorization.k8s.io
```

- `Create yaml specs to launch VSM as K8s deployments & K8s service`

```bash
$ cat my-jiva-vsm.yaml
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: my-jiva-vsm
```

```bash
curl -k -H "Content-Type: application/yaml" \
  -XPOST -d"$(cat my-jiva-vsm.yaml)" \
  http://10.44.0.1:5656/latest/vsm/

{
  "metadata": {
    "annotations": {
      "be.jiva.volume.openebs.io\/count": "2",
      "be.jiva.volume.openebs.io\/vol-size": "1G",
      "iqn": "iqn.2016-09.com.openebs.jiva:my-jiva-vsm",
      "targetportal": "10.96.17.42:3260"
    },
    "creationTimestamp": null,
    "name": "my-jiva-vsm"
  },
  "spec": {
    "AccessModes": null,
    "Capacity": null,
    "ClaimRef": null,
    "OpenEBS": {
      "volumeID": ""
    },
    "PersistentVolumeReclaimPolicy": "",
    "StorageClassName": ""
  },
  "status": {
    "Message": "",
    "Phase": "",
    "Reason": ""
  }
}
```

- `Read an existing VSM`

```bash
curl http://10.44.0.1:5656/latest/vsm/read/<vsm-name>

# e.g.

$ curl http://10.44.0.1:5656/latest/vsm/read/my-jiva-vsm
{
  "metadata": {
    "annotations": {
      "be.jiva.volume.openebs.io\/vol-size": "1G",
      "iqn": "iqn.2016-09.com.openebs.jiva:my-jiva-vsm",
      "targetportal": "10.96.17.42:3260",
      "be.jiva.volume.openebs.io\/count": "2"
    },
    "creationTimestamp": null,
    "name": "my-jiva-vsm"
  },
  "spec": {
    "AccessModes": null,
    "Capacity": null,
    "ClaimRef": null,
    "OpenEBS": {
      "volumeID": ""
    },
    "PersistentVolumeReclaimPolicy": "",
    "StorageClassName": ""
  },
  "status": {
    "Message": "",
    "Phase": "",
    "Reason": ""
  }
}
```

- `Verify the launched Deployments & Services`

```bash
ubuntu@kubemaster-01:~$ kubectl get service
NAME                   CLUSTER-IP    EXTERNAL-IP   PORT(S)             AGE
kubernetes             10.96.0.1     <none>        443/TCP             3d
my-jiva-vsm-ctrl-svc   10.96.17.42   <none>        3260/TCP,9501/TCP   1m
ubuntu@kubemaster-01:~$ 

ubuntu@kubemaster-01:~$ kubectl get services/my-jiva-vsm-ctrl-svc -o json
{
    "apiVersion": "v1",
    "kind": "Service",
    "metadata": {
        "creationTimestamp": "2017-06-12T08:54:04Z",
        "labels": {
            "vsm": "my-jiva-vsm"
        },
        "name": "my-jiva-vsm-ctrl-svc",
        "namespace": "default",
        "resourceVersion": "37744",
        "selfLink": "/api/v1/namespaces/default/services/my-jiva-vsm-ctrl-svc",
        "uid": "b00293e7-4f4c-11e7-bdb8-021c6f7dbe9d"
    },
    "spec": {
        "clusterIP": "10.96.17.42",
        "ports": [
            {
                "name": "iscsi",
                "port": 3260,
                "protocol": "TCP",
                "targetPort": 3260
            },
            {
                "name": "api",
                "port": 9501,
                "protocol": "TCP",
                "targetPort": 9501
            }
        ],
        "selector": {
            "vsm": "my-jiva-vsm-ctrl"
        },
        "sessionAffinity": "None",
        "type": "ClusterIP"
    },
    "status": {
        "loadBalancer": {}
    }
}
```

```
ubuntu@kubemaster-01:~$ kubectl get deploy
NAME               DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
my-jiva-vsm-ctrl   1         1         1            1           5m
my-jiva-vsm-rep    2         2         2            2           5m
```

```
ubuntu@kubemaster-01:~$ kubectl get deploy/my-jiva-vsm-ctrl -o json
{
    "apiVersion": "extensions/v1beta1",
    "kind": "Deployment",
    "metadata": {
        "annotations": {
            "deployment.kubernetes.io/revision": "1"
        },
        "creationTimestamp": "2017-06-12T08:54:04Z",
        "generation": 1,
        "labels": {
            "vsm": "my-jiva-vsm"
        },
        "name": "my-jiva-vsm-ctrl",
        "namespace": "default",
        "resourceVersion": "37787",
        "selfLink": "/apis/extensions/v1beta1/namespaces/default/deployments/my-jiva-vsm-ctrl",
        "uid": "aff7baef-4f4c-11e7-bdb8-021c6f7dbe9d"
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "vsm": "my-jiva-vsm"
            }
        },
        "strategy": {
            "rollingUpdate": {
                "maxSurge": 1,
                "maxUnavailable": 1
            },
            "type": "RollingUpdate"
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "vsm": "my-jiva-vsm"
                }
            },
            "spec": {
                "containers": [
                    {
                        "args": [
                            "controller",
                            "--frontend",
                            "gotgt",
                            "my-jiva-vsm"
                        ],
                        "command": [
                            "launch"
                        ],
                        "image": "openebs/jiva:latest",
                        "imagePullPolicy": "Always",
                        "name": "my-jiva-vsm-ctrl-con",
                        "ports": [
                            {
                                "containerPort": 3260,
                                "protocol": "TCP"
                            },
                            {
                                "containerPort": 9501,
                                "protocol": "TCP"
                            }
                        ],
                        "resources": {},
                        "terminationMessagePath": "/dev/termination-log",
                        "terminationMessagePolicy": "File"
                    }
                ],
                "dnsPolicy": "ClusterFirst",
                "restartPolicy": "Always",
                "schedulerName": "default-scheduler",
                "securityContext": {},
                "terminationGracePeriodSeconds": 30
            }
        }
    },
    "status": {
        "availableReplicas": 1,
        "conditions": [
            {
                "lastTransitionTime": "2017-06-12T08:54:04Z",
                "lastUpdateTime": "2017-06-12T08:54:04Z",
                "message": "Deployment has minimum availability.",
                "reason": "MinimumReplicasAvailable",
                "status": "True",
                "type": "Available"
            }
        ],
        "observedGeneration": 1,
        "readyReplicas": 1,
        "replicas": 1,
        "updatedReplicas": 1
    }
}
```

```
ubuntu@kubemaster-01:~$ kubectl get deploy/my-jiva-vsm-rep -o json
{
    "apiVersion": "extensions/v1beta1",
    "kind": "Deployment",
    "metadata": {
        "annotations": {
            "deployment.kubernetes.io/revision": "1"
        },
        "creationTimestamp": "2017-06-12T08:54:04Z",
        "generation": 1,
        "labels": {
            "vsm": "my-jiva-vsm"
        },
        "name": "my-jiva-vsm-rep",
        "namespace": "default",
        "resourceVersion": "37818",
        "selfLink": "/apis/extensions/v1beta1/namespaces/default/deployments/my-jiva-vsm-rep",
        "uid": "b00ef31b-4f4c-11e7-bdb8-021c6f7dbe9d"
    },
    "spec": {
        "replicas": 2,
        "selector": {
            "matchLabels": {
                "vsm": "my-jiva-vsm"
            }
        },
        "strategy": {
            "rollingUpdate": {
                "maxSurge": 1,
                "maxUnavailable": 1
            },
            "type": "RollingUpdate"
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "vsm": "my-jiva-vsm"
                }
            },
            "spec": {
                "containers": [
                    {
                        "args": [
                            "replica",
                            "--frontendIP",
                            "10.96.17.42",
                            "--size",
                            "1G",
                            "/openebs"
                        ],
                        "command": [
                            "launch"
                        ],
                        "image": "openebs/jiva:latest",
                        "imagePullPolicy": "Always",
                        "name": "my-jiva-vsm-rep-con",
                        "ports": [
                            {
                                "containerPort": 9502,
                                "protocol": "TCP"
                            },
                            {
                                "containerPort": 9503,
                                "protocol": "TCP"
                            },
                            {
                                "containerPort": 9504,
                                "protocol": "TCP"
                            }
                        ],
                        "resources": {},
                        "terminationMessagePath": "/dev/termination-log",
                        "terminationMessagePolicy": "File",
                        "volumeMounts": [
                            {
                                "mountPath": "/openebs",
                                "name": "openebs"
                            }
                        ]
                    }
                ],
                "dnsPolicy": "ClusterFirst",
                "restartPolicy": "Always",
                "schedulerName": "default-scheduler",
                "securityContext": {},
                "terminationGracePeriodSeconds": 30,
                "volumes": [
                    {
                        "hostPath": {
                            "path": "/tmp/my-jiva-vsm/openebs"
                        },
                        "name": "openebs"
                    }
                ]
            }
        }
    },
    "status": {
        "availableReplicas": 2,
        "conditions": [
            {
                "lastTransitionTime": "2017-06-12T08:54:20Z",
                "lastUpdateTime": "2017-06-12T08:54:20Z",
                "message": "Deployment has minimum availability.",
                "reason": "MinimumReplicasAvailable",
                "status": "True",
                "type": "Available"
            }
        ],
        "observedGeneration": 1,
        "readyReplicas": 2,
        "replicas": 2,
        "updatedReplicas": 2
    }
}
```
