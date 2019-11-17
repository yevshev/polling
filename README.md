# HTTP Polling Simulation
*Server-sent Event repository can be found [here](https://github.com/yevshev/server-sent)*

## Information
The `client` directory contains the source for the [http-client](https://hub.docker.com/repository/docker/yevshev/http-client) docker image 

The `server` directory contains the source for the [http-server](https://hub.docker.com/repository/docker/yevshev/http-server) docker image

## Create an overlay network
Create a virtual network for all of our containers to communicate with eachother, and name it 'ssenet':
```sh
$ docker network create -d overlay --attachable ssenet
```

## Deploying to Docker Swarm
Deploy the http server containers defined in [servers.yml](https://github.com/yevshev/polling/blob/master/servers.yml), each running our Go http server binary, and name it 'http':

```sh
$ docker stack deploy -c servers.yml http
```
Deploy the http client container defined in [client.yml](https://github.com/yevshev/polling/blob/master/client.yml), running our Go http client binary, and name it 'polling':

```sh
$ docker stack deploy -c client.yml polling
```

## Client logs
View a real-time feed of all the data the client receives from the servers:
```sh
$ docker service logs polling_client -f
```

## Performance Metrics
View a real-time feed of resource utilization and performance for each runing container:
```sh
$ docker stats
```

## Cleanup
Stop and delete the 'http' server stack and 'polling' client stack:
```sh
$ docker stack rm http polling
```