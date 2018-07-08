

Config : 
```

all config is set enviroment

STORAGE_ROOT : storage root dir 
LISTEN_ADDRESS : listen address 
RABBITMQ_SERVER : rabbitmq address
LOG_DIR : log dir
LOG_LEVEL : log level [warn,info,error,debug]

```


Pre
```

rabbit : message exchange for apiServer and dataServer and Service discovery

docker run --hostname my-rabbit9 --name rabbit-mq9  -p 28090:15672 -p 15672:5672 -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=admin123 rabbitmq:3-management


```


ApiServer
```

    Client Connect for User

```

DataServer
```
    Object Saving service

```

Scripts
```
    scripts:  
        apiServer.sh start apiServer
        dataServer.sh start dataServer 
        test.sh  test 

```