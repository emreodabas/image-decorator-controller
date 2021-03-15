# image-decorator-controller
###### Firstly this project structure is created with [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) 
Kubernetes Controller for importing external images to own/local repository  and redeploy workload

- Docker account -> https://hub.docker.com/u/kubermatico
- Demo run on the latest Kind cluster (kind v0.10.0 go1.15.7 linux/amd64)

[![asciicast](https://asciinema.org/a/xnG9oES86DxMzop2hr8rD7zlp.svg)](https://asciinema.org/a/xnG9oES86DxMzop2hr8rD7zlp)

 ### Run Demo
 - `git clone https://github.com/emreodabas/image-decorator-controller` 
 - `cd /image-decorator-controller && go get && /hack/runDemo.sh ` 
 removing sample data 
 - `hack/resetWorkloads.sh`
 check sample data states
 - `hack/checkStates.sh`
### Installation
#### localhost
- `git clone https://github.com/emreodabas/image-decorator-controller` 
- `cd image-decorator-controller && make run ENV=dev` 

#### Kubernetes 
!! kustomize required
- `git clone https://github.com/emreodabas/image-decorator-controller` 
- `cd image-decorator-controller && go get  && make deploy` 

#### Docker publish
- `git clone https://github.com/emreodabas/image-decorator-controller` 
- `cd image-decorator-controller && make docker-publish` 


### Configurations
for development usage add environment variable ENV=dev

| Variable Name      | Description | Sample |
| ----------- | ----------- | ---------- |
| BACKUP_REGISTRY_ADDRESS      | (required)  backup registry address | "kubermatico" "gcr.io/kubermatico"       |
| REQUEUE_DURATION   |  ms based requeue duration  | "5000"         |
| IGNORED_NS | comma seperated namespaces for ignored namespaces | "kube-system,kube-builder-system" 
| USERNAME | (required) registry username | kubermatico 
| PASSWORD |  required if access token is null | -
| ACCESS_TOKEN | required if password is null | - 
| KUBE_CONFIG_PATH | local path of kubeconfig | /home/emreo/.kube/config
| LEADER_ELECTION | boolean default false | true - false

#### Known issues
* Controller throw exception if resource updated while process
    * This reconciles are requeueing and processed after requeue duration
* Configuration file committed for development usage 
   
[kubebuilder]: https://github.com/kubernetes-sigs/kubebuilder[![asciicast](https://asciinema.org/a/xnG9oES86DxMzop2hr8rD7zlp.svg)](https://asciinema.org/a/xnG9oES86DxMzop2hr8rD7zlp)