

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


elasticsearch : save metadata 

docker run -p 9200:9200 -e "http.host=0.0.0.0" -e "transport.host=127.0.0.1" docker.elastic.co/elasticsearch/elasticsearch:5.0.2

http://127.0.0.1:9200/metadata xput -d'
{
	"mappings":{
		"objects":{
			"properties":{
				"name":{
					"type": "string",
					"index" : "not_analyzed"
				},
				"version":{
					 "type": "integer"
				},
				"size":{
					"type": "integer"
				},
				"hash":{
					"type": "string"
				
				}
			}
		}
		
	}
}'


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