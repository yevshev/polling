# HTTP Polling Simulation
*Server-sent Event repository can be found [here](https://github.com/yevshev/server-sent)*

## Information
The `client` directory contains the source for the [sse-client](https://hub.docker.com/repository/docker/yevshev/sse-client) docker image 

The `server` directory contains the source for the [sse-server](https://hub.docker.com/repository/docker/yevshev/sse-server) docker image

## Deploying to Docker Swarm
Deploy 100 containers with each running our go binary, and name it 'test-cluster':

```sh
$ docker service create --replicas 100 --name test-cluster yevshev/emulator-push

# We can provide addtional parameters as needed
```
Additional docker service [commands](https://docs.docker.com/engine/reference/commandline/service/
):
```sh
# List running docker services
$ docker service ls

# Display detailed information on the service
$ docker service inspect test-cluster --pretty

# Scale our service to 1000 containers
$ docker service scale test-cluster=1000

# Remove the service
$ docker service rm test-cluster

# stop all running containers
$ docker stop $(docker ps -a -q)

# Delete all containers
$ docker rm $(docker ps -a -q)
```
