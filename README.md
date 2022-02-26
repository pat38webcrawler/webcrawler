# WEBCRAWLER
# Introduction 
I made some implementation choices (generally commented in the code) as for 
example the root url passed to the server that must be a full url 
http[s]//base_site (http:// is mandatory). 

As asked in the exercise, all found linked that are not under 
the base url are considered as leafs (so they appear in the displayed sitemap, 
but we don't follow these links).

Only simple unit tests for utils package (for demonstration purposes) have been developed 
(for  other packages, some mock should be written in order to develop unit tests) 

The initial version of the webcrawler (***basic*** git tag) was a naive approach based 
on a recursive algorithm to compute sitemap. 
In later versions, in order to improve performance and avoid stack overflow issues, I get rid of 
recursion to compute the sitemap and I also rely on concurrent goroutines to improve performance
## Setup 
Linux based instructions
```bash
~> git clone https://github.com/pat38webcrawler/webcrawler.git

~> cd  webcrawler

# set  GOPATH accordingly 
# check PATH and GOROOT variables
 
# get required packages 
~> go get -t -v github.com/julienschmidt/httprouter
~> go get -t -v github.com/justinas/alice
~> go get -t -v github.com/pkg/errors
~> go get -t -v golang.org/x/net/html
~> go get -v golang.org/x/net/html
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

~> ./wcclient localhost https://<your_site> 
```
## Build a simple docker image 
This simple image will allow to run our webcrawler server in a docker 
container. Its a first step to be able to deploy  our app in a K8s 
environment

**_NOTE:_**
I use a local docker registry to illustrate that examplee but you can tag the built image and push it under any 
registry you have access

```bash
# go to webcrawler dir
~> cd <webcrawler>
docker build -t node1:5000/webcrawler:latest . 
#then you can start our server in a container 
~> docker run -p 8900:8900 -d node1:5000/webcrawler
```
