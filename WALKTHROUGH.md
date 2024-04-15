# Devops Test Walkthrough

## Basic Tasks
0. Coz we need monitoring later on let's add basic library with prom metrics (check ./pkg/faraway-metrics). 
   Also let's add `healthchecks` library (I provide an example here: ./pgk/faraway-healthchecks)
1. create Dockerfile for this application
   - check .faraway/Dockerfile
   ```dockerfile
    # coz we have 1.16 inside go.mod
    FROM golang:1.20 as builder
    WORKDIR /app
    COPY go.mod go.sum ./
    RUN go mod download
    # also add appropriate .dockerignore could be wise :)
    COPY . .
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o farawayweb .

    # Possible to use scratch images, but skip it for now
    FROM alpine:latest  
    RUN apk --no-cache add ca-certificates
    WORKDIR /root/
    COPY --from=builder /app/farawayweb .
    EXPOSE 8080
    CMD ["./farawayweb"]
   ```

2. create docker-compose.yaml to replicate a full running environment
    so that a developer can run the entire application locally without having
    to run any dependencies (i.e. redis) in a separate process.
   - ./docker-compose.yml
   ```yaml
      services:
        web:
          build:
            context: .
            dockerfile: ./.faraway/Dockerfile
          ports:
            - "8080:8080"
          environment:
            - REDIS_ADDR=redis:6379
          depends_on:
            - redis
          networks:
            - faraway-net
  
        redis:
          image: redis:latest
          ports:
            - "6379:6379"
          networks:
            - faraway-net
  
      networks:
        faraway-net:
          driver: bridge
  
   ```
3. Explain how you would monitor this application in production.   
        We can use Prometheus (or VictoriaMetrics) together with Grafana for metric collection and visualization. 
        We've already integrated Prometheus metrics into our application thanks to our custom library ;))))).
        Depending on environment we probably already have Promethes with Service Discovery based on labels or annotations. 
        So below it is just example to intergrate prom to local docker-compose creatred environment :)

        prometheus.yml for local environment could look like:

        ```yaml
        global:
          scrape_interval: 15s

        scrape_configs:
          - job_name: 'faraway-webapi'
            static_configs:
              - targets: ['web:8080'] 

        ```

        For prod environment we probably already had our monitoring stack. So we need to adapt our service to suit the requirements from SRE/devops team. I'll skip it from here.


## Minikube Tasks
1.  prepare local Kubernetes environment (using MiniKube + Helm) to run our application in pod/container.
it should be created a helm chart with resources for deploying application to Kuberenetes. 
store all relevant scripts (kubectl commands etc) in your forked repository.
  
  I've put charts inside `./.faraway/charts` for webapi and redis. 
  It is basically as simple as `helm create <chart>`. 
  I've add some changes to set correct `env` and also provide our path to live and ready checks thanks to our library those handlers should be avaiable at /ready /live


2.  suggest & create minimal local infrastructure to perform functional testing/monitoring of our application pod.
demonstrate monitoring of relevant results & metrics for normal app behavior and failure(s).  
  
  Because we are limited in time (1 hour) I don't plan to provide any prometheus setup here. 
  I've just point some key things to cosider:
  - using library from your `platfrom engineering team` to make uniformly aligned metrics
  - we would use prometheus service discovery. Please align your services (labels & annotations) accordingly
  - `platform-engineering team` should provide extendable library so that developers could work and add business metrics for their services
  - as for monitoring stack overall - let's stick with pretty basic setup (coz we don't have any functional requirements according to the tasks). So basically Prom (VM) + Grafana + AlertManager (consider using Thanos for lts metrics) should cover our needs.
 


### copy paste from terminal to show we have smtn deployed, up & running:
```bash

~/work/dev/repos/github.com/FarawayGG/devops-test-task main*   [k8s] config  
dev ❯ kgp
NAME                              READY   STATUS    RESTARTS   AGE
faraway-webapi-857c5494fb-gld9l   1/1     Running   0          12h
redis-master-0                    1/1     Running   0          12h

~/work/dev/repos/github.com/FarawayGG/devops-test-task main*   [k8s] config  
dev ❯ helm list                              
NAME            NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                   APP VERSION
faraway-webapi  default         3               2024-04-14 21:28:15.117217907 +0300 MSK deployed        faraway-webapi-0.1.0    1.16.0     
redis           default         1               2024-04-14 21:26:31.707398535 +0300 MSK deployed        redis-19.1.0            7.2.4      

~/work/dev/repos/github.com/FarawayGG/devops-test-task main*   [k8s] config  
dev ❯ k get ing                        
NAME             CLASS   HOSTS                  ADDRESS        PORTS   AGE
faraway-webapi   nginx   faraway-webapi.local   192.168.49.2   80      12h

~/work/dev/repos/github.com/FarawayGG/devops-test-task main*   [k8s] config  
dev ❯ curl faraway-webapi.local/live                                                       
OK
~/work/dev/repos/github.com/FarawayGG/devops-test-task main*   [k8s] config  
dev ❯ curl faraway-webapi.local/    
hello world: updated_time=2024-04-15 06:34:12

~/work/dev/repos/github.com/FarawayGG/devops-test-task main*   [k8s] config  
dev ❯ curl faraway-webapi.local/
hello world: updated_time=2024-04-15 06:34:12

~/work/dev/repos/github.com/FarawayGG/devops-test-task main*   [k8s] config  
dev ❯ curl faraway-webapi.local/
hello world: updated_time=2024-04-15 06:34:12

~/work/dev/repos/github.com/FarawayGG/devops-test-task main*   [k8s] config  
dev ❯ curl faraway-webapi.local/
hello world: updated_time=2024-04-15 06:34:25

~/work/dev/repos/github.com/FarawayGG/devops-test-task main*   [k8s] config  
dev ❯ curl faraway-webapi.local/metrics
# HELP faraway_go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE faraway_go_gc_duration_seconds summary
faraway_go_gc_duration_seconds{quantile="0"} 0
faraway_go_gc_duration_seconds{quantile="0.25"} 0
faraway_go_gc_duration_seconds{quantile="0.5"} 0
faraway_go_gc_duration_seconds{quantile="0.75"} 0
faraway_go_gc_duration_seconds{quantile="1"} 0
faraway_go_gc_duration_seconds_sum 0
faraway_go_gc_duration_seconds_count 0
# HELP faraway_go_goroutines Number of goroutines that currently exist.
# TYPE faraway_go_goroutines gauge
faraway_go_goroutines 14
# HELP faraway_go_info Information about the Go environment.
# TYPE faraway_go_info gauge
faraway_go_info{version="go1.20.14"} 1
# HELP faraway_go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE faraway_go_memstats_alloc_bytes gauge
faraway_go_memstats_alloc_bytes 1.744272e+06
# HELP faraway_go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE faraway_go_memstats_alloc_bytes_total counter
faraway_go_memstats_alloc_bytes_total 1.744272e+06

... 

omitted

...

```

