# dinit

dinit is a simple process supervisor and init system designed to run as PID 1
inside container environments.

All dinit does is manage the startup sequence of the processes in the container, and
wait for it to exit and signal forwarding.


## Features

- startup sequence of the processes
- restart in-place without exiting the container

## Environments

- `DINIT_PORT` - Port used to receive HTTP calls from the remote, default: `8080`
- `DINIT_CONFIG ` - The Configuration of the dinit
- `DINIT_GRACEFUL_TIMEOUT` - Timeout for graceful exit, default: `30`
- `DINIT_EXIT` - Whether the dinit exit when receive signal, default: `true`
- `DINIT_LOG_FILE` - The log file to write to, default: `/data/logs/dinit/dinit.log`
- `DINIT_LOG_FIELDS` - The fileds in log, default: `APP_NAME,POD_IP,HOST_IP,NODE_NAME`

## Configuration file

Example:

```yaml
version: 1.0.0

log_level: debug

services:
  main:
    exec_start: java -server -Xms$(JVM_XMS) -Xmx$(JVM_XMX)...
    lifecycle:
      pre_start:
        exec:
          command: echo "pre_start" && sleep 3
      post_start:
        exec:
          command: echo "post_start" && sleep 3
      pre_stop:
        exec:
          command: echo "pre_stop" && sleep 3
      post_stop:
        exec:
          command: echo "post_stop" && sleep 3
    environment:
      ROUTER: benchmark
    depends_on:
      - consul
  consul:
    exec_start: consul agent -retry-join=consul -config-file /data/consul/config.json
    lifecycle:
      post_start:
        exec:
          command: while true; do sleep 1 && consul catalog services > /dev/null && exit 0 ;done
```

## How to Use

### Docker

- application dockerfile

```
FROM harbor.example.com/public/runtime:slim

WORKDIR /data/apps

ENTRYPOINT ["/urs/local/bin/dinit"]

# java application
CDM ["java", "-server", "-Xms$(JVM_XMS)", "-Xmx$(JVM_XMX)"...]

# python application
CMD ["uwsgi", "--module=runserver", "--callable=app", "--master"...]
```

### Kubernetes

- k8s deployment

```
...
Name:                   xxx-service
Namespace:              default
Labels:                 app=xxx-service
Pod Template:
   dinit:
    Image:      harbor.example.com/public/dinit:1.0.0
    Port:       <none>
    Host Port:  <none>
    Command:
      sh
      -c
    Args:
      echo -en '\n'; cp /usr/local/bin/dinit /data; ls -lh /data/dinit;
    Environment:  <none>
    Mounts:
      /data from dinit (rw)
   xxx-service:
    Image:      harbor.example.com/app/xxx-service:xxxx
    Port:       <none>
    Host Port:  <none>
    Command:
      sh
      -c
    Args:
      echo -en '\n'; cp /data/apps/springboot.jar /data; ls -lh /data/springboot.jar;
    Environment:  <none>
    Mounts:
      /data from dinit (rw)
  Containers:
   app:
    Image:      harbor.example.cn/env/jdk:1.8.0_251
    Port:       8080/TCP
    Host Port:  0/TCP
    Command:
      /data/dinit
    Args:
      java
      -server
      -Xms$(JVM_XMS)
      -Xmx$(JVM_XMX)
      -javaagent:/data/pinpoint-agent/pinpoint-bootstrap.jar
      -Dpinpoint.agentId=$(POD_IP)
      -Dpinpoint.applicationName=$(APP_NAME)
      -Dfile.encoding=UTF-8
      -Djava.awt.headless=true
      -Djava.security.egd=file:/dev/./urandom
      -Xss256k -XX:+AggressiveOpts -XX:-UseBiasedLocking -XX:MaxTenuringThreshold=15 -XX:LargePageSizeInBytes=128m -XX:+UseFastAccessorMethods -XX:MaxGCPauseMillis=20 -XX:InitiatingHeapOccupancyPercent=35 -XX:MetaspaceSize=128M
      -XX:+DisableExplicitGC -XX:+UseG1GC -XX:+PrintGCDateStamps -XX:+PrintGCDetails -Xloggc:/data/logs/gc/gc.log -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/data/logs/gc/heapdump.bin
      -XX:+UnlockExperimentalVMOptions -XX:+UseCGroupMemoryLimitForHeap -XX:+UseContainerSupport
      -jar
      /data/springboot.jar
      --spring.profiles.active=$(SPRING_PROFILES_ACTIVE)
      --server.tomcat.accesslog.rotate=false
      --server.tomcat.accesslog.pattern="%t  %a  %A  %r  %s  %D  %b  %I" >> /data/logs/java/console.log 2>&1
    Limits:
      cpu:     4
      memory:  8Gi
    Requests:
      cpu:      1
      memory:   8Gi
    Readiness:  http-get http://:8080/health delay=80s timeout=1s period=5s #success=1 #failure=3
    Environment:
      TZ:                           Asia/Shanghai
      APP_NAME:                     xxx-service
      APP_REGISTER_NAME:            xxx-service
      JVM_XMS:                      4096m
      JVM_XMX:                      4096m
      APP_ENV:                      <set to the key 'APP_ENV' of config map 'application-environment'>                      Optional: false
      CLOUD:                        <set to the key 'CLOUD' of config map 'application-environment'>                        Optional: false
      ROUTE:                        <set to the key 'ROUTE' of config map 'application-environment'>                        Optional: false
      DINIT_GRACEFUL_TIMEOUT:       <set to the key 'DINIT_GRACEFUL_TIMEOUT' of config map 'application-environment'>       Optional: false
      SPRING_PROFILES_ACTIVE:       <set to the key 'SPRING_PROFILES_ACTIVE' of config map 'application-environment'>       Optional: false
      ZK_HOSTS:                     <set to the key 'ZK_HOSTS' of config map 'application-environment'>                     Optional: false
      CFG_REDIS:                    <set to the key 'CFG_REDIS' of config map 'application-environment'>                    Optional: false
      POD_NAME:                      (v1:metadata.name)
      POD_NAMESPACE:                 (v1:metadata.namespace)
      POD_IP:                        (v1:status.podIP)
      NODE_NAME:                     (v1:spec.nodeName)
      HOST_IP:                       (v1:status.hostIP)
    Mounts:
      ...
   filebeat:
    Image:      harbor.example.cn/public/filebeat:latest
    Port:       <none>
    Host Port:  <none>
    Limits:
      cpu:     1
      memory:  512Mi
    Requests:
      cpu:     100m
      memory:  128Mi
    Environment:
      TZ:  Asia/Shanghai
    Mounts:
      /data/logs from xxx-service (rw)
      /etc/localtime from time (rw)
  Volumes:
    ...
```


### Service Restart In-place

#### Descriptions

The main process restarts these child processes after receiving the `restart In-place` signal


Provide two methods

one example:
```bash
kill -s SIGHUP $dinit_pid
```

another example:

```bash
curl http://localhost:$dinit_port/init/reload
```
