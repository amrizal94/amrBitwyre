# Get Started on Local

## Prerequisite
###  1. You should have docker/docker desktop if on mac/windows
if you not yet installed docker you can download and instal on dockker website on this link:
```
https://www.docker.com/products/docker-desktop/
```
if it's finished, then for mac os download redpanda via brew and rpk for management, you can run this script on your terminal:
```
brew install redpanda-data/tap/redpanda && rpk container start
```
or more information go to github redpanda on this link:
```
https://github.com/redpanda-data/redpanda/
```
### 2. Create a topic on redpanda
- You can create a new topic use rpk which we have downloaded and installed before
- Make sure the container is running, if not yet running you can run this script:
```
rpk container start
```
then create you topic:
```
rpk topic create my-topic
```
change 'my-topic' with the name of the topic you want to create

## Run App
just run this on your path of clonning this repo:
```
go run main.go
```