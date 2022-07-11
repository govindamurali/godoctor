# godoctor

Simple library to parallely ping dependent services and ensure the health. Default implementations for redis, RabbitMq and postgres included. 

**** Example
```
  timeOut:= 2* time.Minute
  handlerFunc:= GetHandler(timeout, RedisChecker(&redisClient), RabbitmqChecker(userName, pasword, host, port), and more...)
```

