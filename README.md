# Technical Test

This repository created for technical test at PT Tiket Digital Raya.

## Installation

Clone the project at go/src

```bash
  git clone https://link-to-project
```

Go to the project directory

```bash
  cd tiket_test
```

use make file 
```bash
make dockerbuild
```

i was including config.json on each service so its easier to deploy.

if you want to change ``port`` for each service there was config.json, please change under 
server object i.e

```
  "server": {
        "address": "0.0.0.0",
        "port": "9922" <- change here
    },

```

## Service server
use postman collection inside this repository and make or replace environment name `server` with value
``
http://172.17.0.1:9922
``
because we using postman at testing web socket we can use postman internal variable also. On **POST** /api/flights/search i'm returning full redis key with format **i.e** ``123123123-1`` you should use full key to access api  **GET** ``/api/flights/search/123123123-1/stream``

## Service Provider
No open rest api or direct accessable api for provider but its running on ``http://172.17.0.1:9911``
## Gravana
You can access gravana at ``http://172.17.0.1:3000``

## Prometheus
You can access prometheus at ``http://172.17.0.1:9090``
## Adminer
i'm simulating server response using delay and postgresql, everything getting delay around 20-30 sec per request to make sure behaviour was right. You can access adminer to access data inside database at ``172.17.0.1:9696`` 
## Postgresql
as for postgre i was deploy it at ``172.17.0.1:5433``
with default
 ```
username=dev_admin
password=dev_password
```
## Architecture
I'm using Custom Clean-Code architecture, using my own boilerplate. As for Provider service, i decided not to include any rest-api because no requirement from  diagrams. As for testing i'm using manual testing because the time kinda short and i don't use any LLM for my codes because i wanted to challenge myself can i implement service without any LLM help. I'm using hard-coded time.Delay() for delaying any input inside provider to simulate network latency and minimize user config-error. The trade off : User can't try at different networking delay.

Contact me if you need help about deployment. Common problem occur usually at server/provider folder try to run go mod tidy to tidy up any package at server/provider folder