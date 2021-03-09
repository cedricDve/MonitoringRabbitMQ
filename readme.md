# First: startup a docker network and run the container
# create docker network
--> docker network create nameNetwork
# run container
docker run -d --rm --net nameNetwork --hostname hostName --name containerName rabbitmq:3.8
# check it !
docker ps 


# ** We can go in the cli with rabbitmqctl
    docker exec -it rabbit-1 bash
        rabbitmqctl
        rabbitmq-plugins => management int.  and prometheus

 # delete a container 
    docker rm -f nameContainer

# run container and browse with management plugin
    docker run -d --rm --net rabbits -p 8080:15672 --hostname rabbit-1 --name rabbit-1 rabbitmq:3.8-management


# create queu & app to add and take messages from the queu 

# publisher
# dockerfile
we use go lang => golang:1.4-alpine 
we need git => dependency manager


# build it, for example build publisher
inside publisher => 
docker build . -t aimvector/rabbitmq-publisher:v1.0.0
# check it 
docker ps

# run publisher app

docker run -it --rm --net rabbits -e RABBIT_HOST=rabbit-a -e RABBIT_PORT=5672 -e RABBIT_USERNAME=guest -e RABBIT_PASSWORD=guest -p 80:80 aimvector/rabbitmq-publisher:v1.0.0
# check it
docker ps
# use post man to check if it work
post a message: localhost:80/publish/messageYouwannaSend


# launch consumer
> inside consumer
>docker build . -t aimvector/rabbitmq-consumer:v1.0.0

>docker run -it --rm --net rabbits -e RABBIT_HOST=rabbit-1 -e RABBIT_PORT=5672 -e RABBIT_USERNAME=guest -e RABBIT_PASSWORD=guest -p 80:80 aimvector/rabbitmq-consumer:v1.0.0


# authentcation
erlang cookie
> grab erlang cookie : 
>> docker exec -it rabbit-a cat /var/lib/rabbitmq/.erlang.cookie
>>> TPLHRMYLVNQVWWVTXYTC

# Cluster formation
# manualy setup a cluster
create two instance with the same Erlang cookie !
> docker run -d --rm --net rabbits --hostname rabbit-a --name rabbit-a -p 8081:15672 -e RABBITMQ_ERLANG_COOKIE=TPLHRMYLVNQVWWVTXYTC rabbitmq:3.8-management
> docker run -d --rm --net rabbits --hostname rabbit-b --name rabbit-b -p 8082:15672 -e RABBITMQ_ERLANG_COOKIE=TPLHRMYLVNQVWWVTXYTC rabbitmq:3.8-management

docker run -d --rm --net rabbits --hostname rabbit-d --name rabbit-d -p 8083:15672 -e RABBITMQ_SERVER_ADDITIONAL_ERL_ARGS="-setcookie TPLHRMYLVNQVWWVTXYTC"  rabbitmq:3.8-management


>> Join a node
> docker exec -it rabbit-b rabbitmqctl stop_app
> docker exec -it rabbit-b rabbitmqctl reset
> docker exec -it rabbit-b rabbitmqctl join_cluster rabbit@rabbit-a
> docker exec -it rabbit-b rabbitmqctl start_app
>> check it 
> docker exec -it rabbit-b rabbitmqctl cluster_status
# cluster status command
>> docker exec -it rabbit-a rabbitmqctl cluster_status
# Automated Clustering
>using config file !

docker run -d --rm --net rabbits `
-v ${PWD}/config/rabbit-a/:/config/ `
-e RABBITMQ_CONFIG_FILE=/config/rabbitmq `
-e RABBITMQ_ERLANG_COOKIE=TPLHRMYLVNQVWWVTXYTC `
--hostname rabbit-a `
--name rabbit-a `
-p 8082:15672 `
rabbitmq:3.8-management

docker run -d --rm --net rabbits `
-v ${PWD}/config/rabbit-b/:/config/ `
-e RABBITMQ_CONFIG_FILE=/config/rabbitmq `
-e RABBITMQ_ERLANG_COOKIE=TPLHRMYLVNQVWWVTXYTC `
--hostname rabbit-b `
--name rabbit-b `
-p 8083:15672 `
rabbitmq:3.8-management


# automatic syncronisation mirroring

# enable federation plugin
docker exec -it rabbit-1 rabbitmq-plugins enable rabbitmq_federation 

> inside the container 
>> docker exec -it rabbit-1 bash

rabbitmqctl set_policy ha-fed \
    ".*" '{"federation-upstream-set":"all", "ha-sync-mode":"automatic", "ha-mode":"nodes", "ha-params":["rabbit@rabbit-a","rabbit@rabbit-b"]}' \
    --priority 1 \
    --apply-to queues