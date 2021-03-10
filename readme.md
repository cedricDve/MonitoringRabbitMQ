# First start a docker network & run the container
## Create docker network
``` docker network create nameNetwork ```
## Run container
``` docker run -d --rm --net nameNetwork --hostname hostName --name containerName rabbitmq:3.8 ```
## Verify 
``` docker ps ```

>>> We can go in the cli with rabbitmqctl    ``` docker exec -it rabbit-1 bash ```
     >>> rabbitmqctl
     >>> rabbitmq-plugins => management int.  and prometheus 

## delete a container 
``` docker rm -f nameContainer ```

# Run container and browse with management plugin enabled
``` docker run -d --rm --net rabbits -p 8080:15672 --hostname rabbit-1 --name rabbit-1 rabbitmq:3.8-management ```

# Create Queue & Run App(Publisher & Consumer)  
Able to publish messages to the queue and take messages from the queue we need to run the publisher- and consumer-app
## Publisher
### Dockerfile
we use go lang 
> golang:1.4-alpine 
we need git 
> git is the dependency manager of go 

### build image
Build the publisher
> inside the publisher ('\rabbitmq\app\publisher')
``` docker build . -t aimvector/rabbitmq-publisher:v1.0.0 ```
>> Verify with docker ps 

### run publisher app
``` docker run -it --rm --net rabbits -e RABBIT_HOST=rabbit-a -e RABBIT_PORT=5672 -e RABBIT_USERNAME=guest -e RABBIT_PASSWORD=guest -p 80:80 aimvector/rabbitmq-publisher:v1.0.0 ```
>> Verify with docker ps

# Use PostMan to verify if it works
> POST a message to the queue: 
```localhost:80/publish/messageYouwannaSend```

## Consumer
> inside the consumer ('\rabbitmq\app\consumer')
``` docker build . -t aimvector/rabbitmq-consumer:v1.0.0 ```

### Run consumer app
``` docker run -it --rm --net rabbits -e RABBIT_HOST=rabbit-1 -e RABBIT_PORT=5672 -e RABBIT_USERNAME=guest -e RABBIT_PASSWORD=guest -p 80:80 aimvector/rabbitmq-consumer:v1.0.0 ```


# Authentcation
Erlang Cookie
> grab Erlang cookie : 
``` docker exec -it rabbit-a cat /var/lib/rabbitmq/.erlang.cookie ```
>>>response: TPLHRMYLVNQVWWVTXYTC

# Cluster formation
## Manualy setup a cluster
### Create two instance with the same Erlang cookie !
``` docker run -d --rm --net rabbits --hostname rabbit-a --name rabbit-a -p 8081:15672 -e RABBITMQ_ERLANG_COOKIE=TPLHRMYLVNQVWWVTXYTC rabbitmq:3.8-management ```
``` docker run -d --rm --net rabbits --hostname rabbit-b --name rabbit-b -p 8082:15672 -e RABBITMQ_ERLANG_COOKIE=TPLHRMYLVNQVWWVTXYTC rabbitmq:3.8-management ``` 


### Join a node
``` docker exec -it rabbit-b rabbitmqctl stop_app```
``` docker exec -it rabbit-b rabbitmqctl reset ```
``` docker exec -it rabbit-b rabbitmqctl join_cluster rabbit@rabbit-a ```
``` docker exec -it rabbit-b rabbitmqctl start_app ```
>> Verify it 
``` docker exec -it rabbit-b rabbitmqctl cluster_status```

# Cluster_status command
``` docker exec -it rabbit-a rabbitmqctl cluster_status```

# Automated Clustering
> Using config file !
## Run container with automated Cluster
```
docker run -d --rm --net rabbits `
-v ${PWD}/config/rabbit-a/:/config/ `
-e RABBITMQ_CONFIG_FILE=/config/rabbitmq `
-e RABBITMQ_ERLANG_COOKIE=TPLHRMYLVNQVWWVTXYTC `
--hostname rabbit-a `
--name rabbit-a `
-p 8082:15672 `
rabbitmq:3.8-management
```


# Mirroring
## Auto. Sync. 

> enable federation plugin
```docker exec -it rabbit-1 rabbitmq-plugins enable rabbitmq_federation ```

> inside the container ( ``` docker exec -it rabbit-1 bash ```)

>> ```
rabbitmqctl set_policy ha-fed \
    ".*" '{"federation-upstream-set":"all", "ha-sync-mode":"automatic", "ha-mode":"nodes", "ha-params":["rabbit@rabbit-a","rabbit@rabbit-b"]}' \
    --priority 1 \
    --apply-to queues
    ```