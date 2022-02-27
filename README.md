# WEBCRAWLER
# Introduction 
I made some implementation choices (generally commented in the code),  for 
example the root url passed to the server  must be a full url 
http[s]//base_site (e.g. http[s]:// is mandatory). 

As asked in the intructions, all found linked that are not under 
the base url are considered as leafs (so they appear in the displayed sitemap, 
but we don't follow these links). It obviously avoids scanning the entire internet :) 

Only simple unit tests for utils package (for demonstration purpose) have been developed 
(for  other packages, some mock should be written in order to develop relevant unit tests) 

The initial version of the webcrawler (***basic*** git tag) was a naive approach based 
on a recursive algorithm to compute sitemap. 
In later versions, in order to improve performance and avoid stack overflow issues, I get rid of 
recursion to compute the sitemap and I also rely on concurrent goroutines to improve performance
## Setup 
Linux based instructions
```bash
~> git clone https://github.com/pat38webcrawler/webcrawler.git

~> cd webcrawler

# set  GOPATH accordingly 
# check PATH and GOROOT variables
 
# get required packages 
~> go get -t -v github.com/julienschmidt/httprouter
~> go get -t -v github.com/justinas/alice
~> go get -t -v github.com/pkg/errors
~> go get -t -v golang.org/x/net/html
```
## Build client and server 

In webcrawler directory 
```bash
~> go build -o webcrawler 
~> cd client 
~> go build -o wcclient 
~> cd ..
```

## Launch server and lauch a request using client 
In a first terminal 
```bash
~> cd <webcrawler dir> 
# start the webcrawler server. It will start the server on localhost. Server listen on port 8900
~> ./webcrawler  
2022/02/26 7:56:59 Starting HTTPServer on address [::]:8900   
```
In another terminal go to webcrawler/client directory
```bash 
~> cd webcrawler/client
# launch the client to send a request to the server

~> ./wcclient localhost 8900 https://<your_site> 
```
## Build a simple docker image 
This simple image will allow to run our webcrawler server in a docker 
container. This is a first step towards deploying our app in a K8s 
environment

> **Note**:
I use a local private docker registry (node1:5000) to illustrate this example, but you can tag the built image as you 
> want in order to push it under any registry you have access.

The **Dockerfile** is available under webcrawler directory 
```bash
# go to webcrawler dir
~> cd <webcrawler>
~> docker build -t node1:5000/webcrawler:latest . 
~> docker push node1:5000/webcrawler:latest 

# then you can start our server in a container 
~> docker run -p 8900:8900 -d node1:5000/webcrawler
```
## Deployment  on a K8S cluster (tested on a k3s distribution)
> **Note**: I initially planned to use microk8s distribution, but I faced some limitations 
> to install it on a WSL2 environment on my laptop (snap not available...). I investigated these issues and made some 
> progress. After some workarounds I was able to use snap on WSL2, but still the microk8s installation failed (controller 
> didn't start). Due to lack of time I decided to rely on k3s distribution which is simpler to install.
> 
>I used latest stable k3s distrib: "k3s version v1.22.6+k3s1"

The **webcrawler.yaml** file contains the different resources needed to deploy my webcrawler container on the k3s cluster.
Basically I created  a service to expose the webcrawler app running in a pod. I also configured an ingress to be 
able to send requests from outside. In this simple example, I only deploy 1 replica for the webcrawler pod, but number of 
replica could be increased (for high availability purpose or to scale the configuration to handle more requests)

The content of the deployment yaml file is 
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: webcrawler-deployment
  namespace: webcrawler-demo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: webcrawler
  template:
    metadata:
      labels:
        app: webcrawler
    spec:
      containers:
        - name: webcrawler
          image: patduc/demo:webcrawler
          resources:
            requests:
              memory: "64Mi"
              cpu: "100m"
            limits:
              memory: "128Mi"
              cpu: "500m"
          ports:
          - containerPort: 8900
          imagePullPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  name: webcrawler-service
  namespace: webcrawler-demo
spec:
  ports:
  - port: 80
    targetPort: 8900
    name: tcp
  selector:
    app: webcrawler
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: webcrawler-ingress
  namespace: webcrawler-demo
  annotations:
    kubernetes.io/ingress.class: "traefik"
spec:
  rules:
  - http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: webcrawler-service
            port:
              number: 80
```
> I have preliminarily pushed my webcrawler docker image in a dockerhub repository ***patduc/demo:webcrawler***
I also wanted to try with a local unsecured registry by modifying the rancher ***registries.yaml*** file, but 
I faced some issues (k3s try to access the registry as a secured one though it is unsecured. It was probably a 
configuration mistake on my side, but I didn't push the investigation deeper to save time)

In order to deploy the different elements I ran the following commands 

```bash
# Create a dedicaated namespace
~> kubectl create namespace webcrawler-demo

#  deploy service and webcrawler pod
~> k3s kubectl apply -f webcrawler.yaml
deployment.apps/webcrawler-deployment created
service/webcrawler-service created
ingress.networking.k8s.io/webcrawler-ingress created

# check service and pod are running 
~> k3s kubectl get ingress,svc,pods -n webcrawler-demo
NAME                                           CLASS    HOSTS   ADDRESS      PORTS   AGE
ingress.networking.k8s.io/webcrawler-ingress   <none>   *       10.0.0.161   80      32m26s

NAME                         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
service/webcrawler-service   ClusterIP   10.43.240.114   <none>        80/TCP    32m26s

NAME                                         READY   STATUS    RESTARTS   AGE
pod/webcrawler-deployment-64f74bd9bb-m2cq9   1/1     Running   0          32m26s
```

Finally, I was able to  run a request using our client

```bash
~> ./wcclient 10.197.138.254  80  https://access.redhat.com/products/
https://access.redhat.com/products/
  - https://access.redhat.com/products/#pfe-navigation
    - https://access.redhat.com/products/#cp-main
    - https://access.redhat.com/management/
    - https://access.redhat.com/downloads/
    - https://catalog.redhat.com/software/containers/explore/
    - https://access.redhat.com/support/cases/
    - https://access.redhat.com/
    ...
```
