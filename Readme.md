# kube-query
[![LICENSE](https://img.shields.io/github/license/Shadow-linux/kube-query
)](https://github.com/Shadow-linux/kube-query/blob/master/LICENSE)

A kubectl plug-in that makes it easier to query and manipulate K8S clusters.
[(what is kubectl plug-in ?)](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)

![demo](./daemon.gif)

Kube-query support some resource shortcut query, like: pods, deploy, service, configmap, dameonset, job, nodes. We can get their simple information, relationship or login container.

Kube-query accepts the same commands as the kubectl, except you don't need to provide the kubectl prefix. So it doesn't require the additional cost to use this cli.

## Installation

#### Downloading standalone binary
Binaries are available from [(github release)](https://github.com/Shadow-linux/kube-query/releases).

<details>
<summary>macOS (darwin) - amd64</summary>

```
wget https://github.com/Shadow-linux/kube-query/releases/download/v1.0.0/kube-query_v1.0.0_darwin_amd64.zip
unzip kube-query_v1.0.0_darwin_amd64.zip
chmod +x kube-query.darwin-amd64
sudo mv bin/kube-query.darwin-amd64 /usr/local/bin/kube-query
```

</details>

<details>
<summary>Linux - amd64</summary>

```
wget https://github.com/Shadow-linux/kube-query/releases/download/v1.0.0/kube-query_v1.0.0_linux_amd64.zip
unzip kube-query_v1.0.0_linux_amd64.zip
chmod +x kube-query.linux-amd64
sudo mv bin/kube-query.linux-amd64 /usr/local/bin/kube-query
```

</details>

<details>
<summary>Source code</summary>

```
# install go version 1.17+ first.
wget https://github.com/Shadow-linux/kube-query.git
cd kube-query
make build
mv bin/kube-query /usr/local/bin/kube-query
```

</details>




## Goal

Hopeful easier to query and manipulate K8S cluster.

## Usage
* config your kubeconfig
```shell
export KUBECONFIG=~/.kube/config
```

#### Start way
1. use in kubectl
```shell
mv /usr/local/bin/kube-query /usr/local/bin/kubectl-query
kubectl query [--debug]
```

2. standalone
```shell
./kube-query [--debug]
```

#### Basic Command

* Clear console
```shell
kube-query ~ > clear
```

* Show help info for `kube-query` and `kubectl` native command.
```shell
kube-query ~ > help
```

* Set global namespace
```shell
kube-query ~ > use default
Set namespace default
kube-query ~ > use kube-system
Set namespace kube-system
# all namespace can not use in native kubectl command.
# set namespace all is not safety operation. 
kube-query ~ > use all
Set namespace all
```

* Run `shell` command. We can easier auto complete the file path, when you input a word start with  `/` or `./`.
```shell
kube-query ~ > @ ls /tmp;
```

* Close console.
```shell
# exit | quit | ctrl + D
kube-query ~ > exit
```
---
#### Resource command
> resource format: ResourceName.Namespace,
* Output mode.
```shell
kube-query ~ > pods jtthink-ngx-8669b5c9d-xwljg.default [-o desc]
kube-query ~ > pods jtthink-ngx-8669b5c9d-xwljg.default -o [yaml|desc|json]
```

* Show labels
```shell
kube-query ~ > pods jtthink-ngx-8669b5c9d-xwljg.default -l
Labels: 
app=jtthink-ngx
pod-template-hash=8669b5c9d
```

* Show events
```shell
kube-query ~ > pods jtthink-ngx-8669b5c9d-xwljg.default -e
Events: 
TYPE    REASON  MESSAGE 
```

* Show relationship
```shell
kube-query ~ > pods jtthink-ngx-8669b5c9d-xwljg.default -r
Relevant relationship:
##### Service #####
NAME                    TYPE            CLUSTER-IP      EXTERNAL-IP     PORTS        
jtthink-ngx-svc         ClusterIP       10.99.226.31                    38080/TCP       
jtthink-ngx-svc-1       NodePort        10.99.47.202                    81:30080/TCP    

##### ReplicaSet #####
NAME                    DESIRED CURRENT READY 
jtthink-ngx-8669b5c9d   1       1       1       

##### Deployment #####
NAME            READY   UP-TO-DATE      AVAILABLE 
jtthink-ngx     1/1     1               1   
```
* Connect to container
```shell
# connect specify container name and shell.
kube-query ~ > pods jtthink-ngx-8669b5c9d-xwljg.default -i jt-nginx -s /bin/sh
Connect to container: jtthink-ngx-8669b5c9d-xwljg.jt-nginx
Commannd: /bin/sh 
/ # exit
Connection closed.

# connect default container.
kube-query ~ > pods jtthink-ngx-8669b5c9d-xwljg.default -i
Connect to container: jtthink-ngx-8669b5c9d-xwljg.
Commannd: sh 
/ # exit
Connection closed.

```

* Use `grep` command to filter info;
```shell
kube-query ~ > pods jtthink-ngx-8669b5c9d-xwljg.default -l |grep -A 1  -i labels
Labels: 
app=jtthink-ngx
```

#### Native kubectl command

* example: get
```shell
# we can use `Tab` to auto complete. 
kube-query ~ > get nodes k8s-01;
kube-query ~ > get pods;
kube-query ~ > get svc;
```

The same way you normally use Kubectl anyway, just you do not need input `kubectl`. 


## Author

ShadowYD

* Mail: 972367265@qq.com
* Juejin: [@ShadowYD](https://juejin.cn/user/2524134427859960)

## LICENSE

This software is licensed under the MIT License (See [LICENSE](./LICENSE)).
