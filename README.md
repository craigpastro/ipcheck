# Anomaly Detection

This is the gRPC version. There is also a [Gin version](https://github.com/siyopao/ipcheck/tree/gin-version).

## Design document

 - [Design document](https://docs.google.com/document/d/1i_hwcNFGmx_v72G_TZ9YjHjzUM6Yv74tvBlvb_CoHfU/edit#) which **is now badly out of date**.

 Let me try to explain some of the differences:
 - So now we make use of an in-memory trie to do the IP address checking.
 - Our response no longer contains which lists that the IP address belongs to. It just returns "yes" or "no".
 - There is an endpoint `PUT /v1/addresses` which is use to tell the system to update its blocklists from the blocklist repo. If it fails to update, the system will just crash. The expectation is that you will spin up a new one.

## Usage

### Install Go

https://golang.org/

### Configuration

Set any variables in [`.env`](./.env) to configure the behaviour of the service.

- `SERVER_ADDR` - The address the server will listen on. The default is `127.0.0.1:8080`.
- `IP_SETS_DIR` - The directory where you want to download the ipsets to. Default is `/tmp/ipsets` (which will not work on Windows so please update it if you are using Windows).
- `IP_SETS` - A comma separated list of blocklists that you would like the system to use. Please find the list of all blocklists at https://github.com/firehol/blocklist-ipsets. Note, adding many blocklists may lead to performance degradation. The default is `feodo.ipset,palevo.ipset,sslbl.ipset,zeus.ipset,zeus_badips,dshield.netset,spamhaus_drop.netset,spamhaus_edrop.netset,fullbogons.netset,openbl.ipset,blocklist_de.ipset`. (The Level 1 and some of Level 2 blocklists described [here](https://github.com/firehol/blocklist-ipsets#which-ones-to-use).)

### Running the service

1. If the gRPC code needs updating, regenerate it with:
  ```
  protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        api/proto/v1/ipcheck.proto
  ```
2. Then
```
go run .
```
should start the service listening at `[::1]:50051`.

### API

We will use [`grpcurl`](https://github.com/fullstorydev/grpcurl) for the examples.

- Check if `193.242.145.0` is in the blocklist with:
  ```
  grpcurl -plaintext -import-path ./api/proto/v1 -proto ipcheck.proto -d '{"ip": "193.242.145.0"}' [::1]:50051 api.proto.v1.IpCheck/InBlocklist
  ```
- Update the blocklists with:
  ```
  grpcurl -plaintext -import-path ./api/proto/v1 -proto ipcheck.proto [::1]:50051 api.proto.v1.IpCheck/InitBlocklists
  ```

### Tests

Run the tests:
```
go test ./...
```

## Benchmarking

```
$ ab -n 10000 -c 10 http://127.0.0.1:8080/v1/addresses/193.242.145.28
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 1000 requests
Completed 2000 requests
Completed 3000 requests
Completed 4000 requests
Completed 5000 requests
Completed 6000 requests
Completed 7000 requests
Completed 8000 requests
Completed 9000 requests
Completed 10000 requests
Finished 10000 requests


Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /v1/addresses/193.242.145.0
Document Length:        20 bytes

Concurrency Level:      10
Time taken for tests:   4.514 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1430000 bytes
HTML transferred:       200000 bytes
Requests per second:    2215.29 [#/sec] (mean)
Time per request:       4.514 [ms] (mean)
Time per request:       0.451 [ms] (mean, across all concurrent requests)
Transfer rate:          309.36 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.6      0      36
Processing:     0    4  29.3      0     305
Waiting:        0    4  29.3      0     305
Total:          0    4  29.3      1     305

Percentage of the requests served within a certain time (ms)
  50%      1
  66%      1
  75%      1
  80%      1
  90%      1
  95%      1
  98%      3
  99%    168
 100%    305 (longest request)
```

## Docker

You can build the image with
```
docker build -t ipcheck .
```
Then can run the container using
```
docker run --rm --name ipcheck -p 127.0.0.1:50051:50051 -e SERVER_ADDR=":50051" ipcheck
```
When that starts listening you can make requests to `127.0.0.1:50051`.

## Steps to make IP Checker production ready

Some steps are described in the [Future Work](https://docs.google.com/document/d/1i_hwcNFGmx_v72G_TZ9YjHjzUM6Yv74tvBlvb_CoHfU/edit#heading=h.bcsw102vr267) section of the design document.
