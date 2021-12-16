clear nsq topic
---


### 1、get topic
request `GET` `/api/topics`
response
```
{"topics":["test"],"message":""}
```

### 2、get topic depth
requets `GET` `/api/topics/<topic>`
response
```
{
    "node":"*",
    "hostname":"",
    "topic_name":"test",
    "depth":0,
    "memory_depth":0,
    "backend_depth":0,
    "message_count":202,
    ...
    "message":""
}
```

### 3、empty queue
request `POST` `/api/topics/<topic>`
body
```
{action: "empty"}
```
