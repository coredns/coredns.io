+++
date = "2022-03-03T08:00:00Z"
description = "A guide for using CoreDNS in Apache APISIX."
tags = ["Apache APISIX", "Service", "Discovery", "API gateway"]
title = "CoreDNS and Apache APISIX open new doors for Service Discovery?"
author = "Zijie Chen"
+++

# Background information

In traditional physical machine and virtual machine deployment, calls between various services can be made through fixed **IP** + **port**. With the advent of the cloud-native era, enterprise business deployment is more inclined to cloud-native containerization. However, in a containerized environment, the startup and destruction of service instances are very frequent. Manual maintenance by operation and maintenance personnel will not only be a heavy workload, but also ineffective. Therefore, a mechanism is needed that can automatically detect the service status, and dynamically bind a new address when the service address changes. The service discovery mechanism came into being.。

# Service Discovery

The service discovery mechanism can be split into two parts:

- Service Registry: Store host and port information for services.

If a container provides a service for calculating the average, we use the service name of `average` as the unique identifier, then it will be stored in the form of a key-value pair (`average:192.168.1.21`) in the service registry.

- Service Discovery: Allows other users to discover the information stored during the service registration phase. It is divided into client discovery mode and server discovery mode.

**Client Service Discovery Mode**

When using the client discovery mode, the client obtains the actual network address of the available service by querying the storage information of the service registry, selects an available service instance through a load balancing algorithm, and sends the request to the service.

Advantages: Simple architecture, flexible expansion, and easy implementation of load balancing functions.

Disadvantages: heavy client, strong coupling, there is a certain development cost.

{{< figure src="/images/apisix-coredns-flowdiagram-1.png" title="Client Discovery Mode.">}}

The implementation logic of client discovery mode is as follows:

1. When a new service is started, it will actively register with the registration center, and the service registration center will store the service name and address of the new service;
2. When the client needs this service, it will use the service name to initiate a query to the service registry;
3. The service registry returns the available addresses, and the client selects one of the addresses to initiate the call according to the specific algorithm.

In this process, in addition to service registration, the work of service discovery is basically completed by the client independently, and the addresses of the registry and the server are also fully visible to the client.

**Server Service Discovery Mode**

The client sends a request to the Load Balancer, and the Load Balancer queries the service registry according to the client's request, finds an available service, and forwards the request to the service. Like the client service discovery mode, the service needs to be registered and deregistered in the registry. 

Advantages: The discovery logic of the service is transparent to the client.

Disadvantages: Requires additional deployment and maintenance of a Load Balancer.

{{< figure src="/images/apisix-coredns-flowdiagram-2.png" title="Server Discovery Mode.">}}

The implementation logic of server discovery mode is as follows:

1. When a new service is started, it will actively register with the registry, and the service registry will store the service name and address of the new service;
2. When the client needs a service, it will use the service name to initiate a query to the load balancer;
3. According to the service name requested by the client, the Load Balancer proxies the client to initiate a request to the service registry;
4. After the Load Balancer obtains the returned address, it selects one of the addresses to initiate the call according to the specific algorithm.

# Advantages of using CoreDNS

Compared with common service discovery frameworks (Zookeeper and Consul), what are the advantages of CoreDNS implementing service discovery?

The principle of service discovery is similar to DNS domain name system, which is an important infrastructure in computer networks. The DNS domain name system binds domain names that rarely change with frequently changing server IP addresses, while the service discovery mechanism is to seldom change domain names. The service name is bound to the service address. In this way, we can use DNS to achieve a function similar to the service registry, and only need to convert the domain name stored in the DNS into the service name. Since many computers have built in DNS functions, we only need to modify the configuration on the original DNS system without doing too many extra things.

CoreDNS is an open source DNS server written in `Go`, which is commonly used for DNS services and service discovery in multi-container environments due to its flexibility and extensibility. CoreDNS is built on top of Caddy, the HTTP/2 web server, and implements a plug-in chain architecture, abstracting many DNS related logic into layer-by-layer plug-ins, which are more flexible and easy to expand, and user selected plugin It will be compiled into the final executable file, and the running efficiency is also very high. CoreDNS is the first cloud native open source project to join CNCF (Cloud Native Computing Foundation) and has graduated, and it is also the default DNS service in Kubernetes.

As middleware, Apache APISIX also integrates a variety of service discovery capabilities. The following will show you how to configure CoreDNS in Apache APISIX.

# Principle Architecture

1. The client initiates a request to APISIX to call the service.
2. APISIX accesses the upstream service node according to the set route (the specific configuration can be seen below). In APISIX, you can set the upstream information to obtain through DNS. As long as the IP address of the DNS server is set correctly, APISIX will automatically initiate a request to this address , to obtain the address of the corresponding service in DNS.
3. CoreDNS returns a list of available addresses based on the requested service name.
4. APISIX selects one of the available addresses and the configured algorithm to initiate the call.

The overall structure is as follows:

{{< figure src="/images/apisix-coredns-flowdiagram-3.png" title="Principle Architecture.">}}

# How to Use

## Prerequisites

This article is based on the following environments.

- OS: Centos 7.9.

- Apache APISIX 2.12.1, please refer to: [How-to-Bulid Apache APISIX](https://apisix.apache.org/docs/apisix/how-to-build).

- CoreDNS 1.9.0，please refer to：[CoreDNS Installation Guide.](https://coredns.io/manual/toc/#installation)

- Node.js, please refer to: [Node.js Installation](https://github.com/nodejs/help/wiki/Installation).

## Procedure

1. Use Node.js's `Koa` framework starts a simple test service on port `3005`。

Accessing this service will return the string `Hello World`, and we will get the address of this service via CoreDNS later.

~~~text
const Koa = require('koa');
const app = new Koa();

app.use(async ctx => {
  ctx.body = 'Hello World';
});

app.listen(3005);
~~~

2. Configure CoreDNS.

CoreDNS listens on port `53` by default, and will read the `Corefile` configuration file in the same directory. Under the initial conditions, there is no `Corefile` file in the same directory, so we need to create and complete the configuration.

`Corefile` mainly implements functions by configuring plugins, so we need to configure three plugins:

- `hosts` ：You can use this parameter to bind the service name and IP address. `fallthrough `means that when the current plugin cannot return normal data, the request can be forwarded to the next plugin for processing (if it exists).
- `forward`：Indicates to proxy the request to the specified address, usually the authoritative DNS server address.
- `log`：Do not configure any parameters to print log information to the console interface for debugging.

~~~text
.:1053 {                           # Listen on port 1053
    hosts {                        
        10.10.10.10 hello    
           # Bind the service name "coredns" to the IP address
        fallthrough
    }
    forward . 114.114.114.114:53  
    log
}
~~~

3. Configuring Apache APISIX。

Add the relevant configuration in the `conf/config.yaml` file and reload Apache APISIX.

~~~text
# config.yml
# ...other config
discovery:   
   dns:     
     servers:       
        - "127.0.0.1:1053"  # Use the real address of the DNS server,
                            # here is the 1053 port of the local machine.
~~~

4. Configure routing information in Apache APISIX.

Next, we configure the relevant routing information by requesting the `Admin API`.

~~~
curl http://127.0.0.1:9080/apisix/admin/routes/1 -H 'X-API-KEY: edd1c9f034335f136f87ad84b625c8f1' -X PUT -d '
{
    "uri": "/core/*",
    "upstream": {
        "service_name": "hello:3005",   
                    # Name the service as coredns, consistent with 
                    # the configuration of the hosts plugin in CoreDNS
        "type": "roundrobin",
        "discovery_type": "dns" # Set service discovery type to DNS
    }
}'
~~~

5. Verify.

a. Authenticate on the local machine

~~~
curl 127.0.0.1:9080/core/hello -i

HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Content-Length: 11
Connection: keep-alive
Date: Wed, 16 Feb 2022 08:44:08 GMT
Server: APISIX/2.12.1

Hello World
~~~

b. Verify on other hosts

~~~
curl 10.10.10.10:9080/core/hello -i

HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Content-Length: 11
Connection: keep-alive
Date: Wed, 16 Feb 2022 08:43:32 GMT
Server: APISIX/2.12.0

Hello World
~~~

As you can see from the above results, the service is running normally.

6. The IP address of the simulated container is changed because the container cannot provide services for various reasons.

We need to set up the same service on another server, also running on port `3005`, but with the IP address changed, and the return string changed to `Hello, Apache APISIX`.

~~~
const Koa = require('koa');
const app = new Koa();

app.use(async ctx => {
  ctx.body = 'Hello, Apache APISIX';
});

app.listen(3005);
~~~

Modify the `Corefile` configuration and restart Core DNS. Leave other configurations unchanged. The configuration example is as follows:

~~~
.:1053 {                           # Listen on port 1053
    hosts {                        
        10.10.10.10 hello    
           # Bind the service name "coredns" to the IP address
        fallthrough
    }
    forward . 114.114.114.114:53  
    log
}
~~~

> DNS has a caching mechanism. When we use the `dig` command to request to resolve a new domain name, we will see a number field in the returned `DNS record`, that is, the `TTL` field, which is generally `3600`, which is one hour. Requests sent to the domain name within the `TTL` period will no longer request the DNS server to resolve the address, but will directly obtain the address corresponding to the domain name in the local cache.

By verifying, we can find that the request has been redirected to the new address. Verify as follows:

~~~
curl 127.0.0.1:9080/core/hello -i

HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Content-Length: 11
Connection: keep-alive
Date: Wed, 16 Feb 2022 08:44:08 GMT
Server: APISIX/2.12.0

Hello, Apache APISIX
~~~

# Summary

This article mainly introduces the types of service discovery and how to use CoreDNS in Apache APISIX. You can use Apache APISIX and CoreDNS according to your business needs and past technical architecture.
