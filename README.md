# Maya API Server

[![GoDoc](https://godoc.org/github.com/openebs/mayaserver?status.png)](https://godoc.org/github.com/openebs/mayaserver)

> OpenEBS exposes its APIs here

A service exposing `Kubernetes` like volume APIs.

## Setting up maya api server development & run environment

> These are some of the steps to start off with development & running of maya api 
server in one's laptop. It assumes use of Linux as laptop's OS. In addition, the 
laptop should have Virtual Box & Vagrant installed.

```bash
- git clone https://github.com/openebs/mayaserver.git
- cd to above cloned folder i.e mayaserver
  - vagrant up
  - vagrant ssh
- Inside the vagrant VM run below steps:
  - make init
  - make
  - make bin
  - sudo nohup m-apiserver up -bind=172.28.128.4 &>mapiserver.log &
```

### Troubleshooting during local setup

```bash
- `make init` is a time taking operation
  - This downloads all the vendoring libraries
  - Typically required for the very first attempt only
  - In case of add/update of new/existing vendoring libraries:
    - use `make sync` than `make init`
```

## Installing maya api server

> Steps to install maya api server's released binary

```bash
- Navigate to https://github.com/openebs/mayaserver/releases
- Download the binary depending on the required release & host OS architecture
  - e.g. Below link points to release `0.0.6` & `linux_386` architecture
  - https://github.com/openebs/mayaserver/releases/download/0.0.6/mayaserver-linux_386.zip
- Extract the m-apiserver binary from above .zip file & put it at /usr/local/bin/
- Set appropriate user groups & permissions
- Start maya api server as a long running daemon
  - sudo nohup m-apiserver up -bind=<IP_Addr> &>mapiserver.log &
```

## Use Cases

### Provisioning OpenEBS Volumes K8s way

Maya api server has some similarity with Kubernetes with respect to latter's
volume provisioning methods. This makes it super simple for operators and admins
to get into the usage of OpenEBS without much of learning curve.

## Maya api server's exposed REST APIs

- `NOTE: Use the bind address on which your maya api server is running`

- Get InstanceID

  ```bash
    # Metadata
    $ curl http://172.28.128.4:5656/latest/meta-data/instance-id
  ```

- Volume provisioning & deletion requires the presence of orchestrator config file
  - It is a .INI formatted file
  - It should be located at `/etc/mayaserver/orchprovider/`
  - There can be multiple .INI files based on the number of `regions` the orchestrator is deployed/running
    - `/etc/mayaserver/orchprovider/nomad_global.INI`
    - `/etc/mayaserver/orchprovider/nomad_us-east-1.INI`
  - Each region based .INI file can be categorised into multiple datacenters
  - Maya api server works with following defaults:
    - `Nomad` as the default orchestrator &
    - `global` as Nomad's default region &
    - `dc1` as Nomad's default datacenter
  - Below is a sample .INI file that is valid for `Nomad` as maya api server's orchestrator

  ```ini
  [datacenter "dc1"]
  ; Address of Nomad deployment
  address = http://172.28.128.3:4646
  
  ; Container Networking options
  cn-type = host
  cn-network-cidr = 172.28.128.1/24
  cn-interface = enp0s8
  
  ; Container Storage options
  cs-persistence-location = /tmp/
  cs-replica-count = 2
  ```

- Below is a sample volume spec that can be provisioned

  ```yaml
  # Similar to K8s' PersistentVolumeClaim
  kind: PersistentVolumeClaim
  apiVersion: v1
  metadata:
    name: minjvol
  spec:
    accessModes:
      - ReadWriteOnce
    resources:
      requests:
        storage: 3G
  ```

  - NOTE - Verbose specs can be found at:
    - `lib/mockit/verbose_openebs_pvc.yaml`
    - json formatted: `lib/mockit/verbose_openebs_pvc.json`
  - Both json & yaml specs are supported when http request's Content-Type is `application/yaml`

- Sample REST Calls
 
  ```bash
  
  # Provision
  
  $ curl -k -H "Content-Type: application/yaml" \
    -XPOST -d"$(cat lib/mockit/minimal_openebs_pvc.yaml)" \
    http://172.28.128.4:5656/latest/volumes/
    
  {
    "Spec": {
      "AccessModes": null,
      "Capacity": null,
      "ClaimRef": null,
      "OpenEBS": {
        "volumeID": ""
      },
      "PersistentVolumeReclaimPolicy": "",
      "StorageClassName": ""
    },
    "Status": {
      "Message": "",
      "Phase": "",
      "Reason": "complete"
    },
    "annotations": {
      "evalpriority": "50",
      "evaltype": "service",
      "evaltrigger": "job-register",
      "evaljob": "myjivavol",
      "evalstatus": "complete",
      "evalstatusdesc": "",
      "evalblockedeval": "360ef0b3-e6dc-9353-d91c-d7c91c239102"
    },
    "creationTimestamp": null,
    "name": "myjivavol"
  }

  # Info
  
  $ curl http://172.28.128.4:5656/latest/volume/info/myjivavol
  
  {
    "Spec": {
      "AccessModes": null,
      "Capacity": null,
      "ClaimRef": null,
      "OpenEBS": {
        "volumeID": ""
      },
      "PersistentVolumeReclaimPolicy": "",
      "StorageClassName": ""
    },
    "Status": {
      "Message": "",
      "Phase": "",
      "Reason": "running"
    },
    "annotations": {
      "be.jiva.volume.openebs.io/vol-size": "3G",
      "fe.jiva.volume.openebs.io/ip": "172.28.128.15",
      "iqn": "iqn.2016-09.com.openebs.jiva:myjivavol",
      "targetportal": "172.28.128.15:3260",
      "JIVA_REP_IP_0": "172.28.128.16",
      "JIVA_REP_IP_1": "172.28.128.8",
      "be.jiva.volume.openebs.io/count": "2"
    },
    "creationTimestamp": null,
    "name": "myjivavol"
  }

  # Delete
  
  $ curl http://172.28.128.4:5656/latest/volume/delete/myjivavol
  
  {
    "Spec": {
      "AccessModes": null,
      "Capacity": null,
      "ClaimRef": null,
      "OpenEBS": {
        "volumeID": ""
      },
      "PersistentVolumeReclaimPolicy": "",
      "StorageClassName": ""
    },
    "Status": {
      "Message": "",
      "Phase": "",
      "Reason": "complete"
    },
    "annotations": {
      "evaltype": "service",
      "evaltrigger": "job-deregister",
      "evaljob": "minjvol",
      "evalstatus": "complete",
      "evalstatusdesc": "",
      "evalblockedeval": "",
      "evalpriority": "50"
    },
    "creationTimestamp": null,
    "name": "myjivavol"
  }

  # Info again
  
  $ curl http://172.28.128.4:5656/latest/volume/info/myjivavol

  Unexpected response code: 404 (job not found)
    
  ```

## Troubleshooting

- Verify the presence of maya api server binary
  - which m-apiserver
  - m-apiserver version

- Verify the presence of maya api server's orchestrator's .INI file(s)
  - i.e. /etc/mayaserver/orchprovider/nomad_global.INI
  - `global` is the name of the region

- Verify the contents of maya api server's orchestrator's .INI file

- Verify if maya api server is running as a process
  - Watch out for the process with 5656 as the port
  - `5656` is the default tcp port on which maya api server's services are exposed

  ```bash
  # Use netstat command
  $ netstat -tnlp

  (Not all processes could be identified, non-owned process info
   will not be shown, you would have to be root to see it all.)
  Active Internet connections (only servers)
  Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
  tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN      -
  tcp        0      0 127.0.0.1:5656          0.0.0.0:*               LISTEN      -
  tcp6       0      0 :::22                   :::*                    LISTEN      -

  # Using sudo will display the PID details
  $ sudo netstat -tnlp

  Active Internet connections (only servers)
  Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
  tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN      1258/sshd
  tcp        0      0 127.0.0.1:5656          0.0.0.0:*               LISTEN      3078/m-apiserver 
  tcp6       0      0 :::22                   :::*                    LISTEN      1258/sshd
  ```

- How to get the iscsi portal & iqn information of a iscsi based jiva volume ?
  - Info based REST API will fetch these information.
  - However, these will **not** be fetched in case of:
    - `error` or
    - `in-progress` / `pending` status

## Licensing

Maya api server is completely open source and bears an Apache license. Maya api 
server's core components and designs are a derivative of other open sourced libraries 
like Nomad and Kubernetes.
