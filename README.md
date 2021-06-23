# Anomaly Detection

## Design document

(**This document is now outdated.**) There is a [design document](https://docs.google.com/document/d/1i_hwcNFGmx_v72G_TZ9YjHjzUM6Yv74tvBlvb_CoHfU/edit#) which describes the general system, assumptions, tradeoffs, future work, etc.

## Usage

### Install Go

https://golang.org/

### Configuration

Set any variables in [`.env`](./.env) to configure the behaviour of the service.

- `SERVER_ADDR` - The address the server will listen on. The default is `127.0.0.1:8080`.
- `GIN_MODE` - The mode Gin will run in. Choose either `debug` or `release`. The default is `release`.
- `IP_SETS_DIR` - The directory where you want to download the ipsets to. Default is `/tmp/ipsets` (which will not work on Windows so please update it if you are using Windows).
- `IP_SETS` - A comma separated list of blocklists that you would like the system to use. Please find the list of all blocklists at https://github.com/firehol/blocklist-ipsets. Note, adding many blocklists may lead to performance degradation. The default is `feodo.ipset,palevo.ipset,sslbl.ipset,zeus.ipset,zeus_badips,dshield.netset,spamhaus_drop.netset,spamhaus_edrop.netset,fullbogons.netset,openbl.ipset,blocklist_de.ipset`. (The Level 1 and some of Level 2 blocklists described [here](https://github.com/firehol/blocklist-ipsets#which-ones-to-use).)

### Running the service

```
go run .
```
should start the service listening at `127.0.0.1:8080`.

### API

Suppose you would like to check if `193.242.145.0` is in the blocklist. Simply
```
curl -XGET http://127.0.0.1:8080/v1/addresses/193.242.145.28
```
The response will be either
```
{
  "inBlocklist": true
}
```
or
```
{
  "inBlocklist": false
}
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

## Steps to make IP Checker production ready

Some steps are described in the [Future Work](https://docs.google.com/document/d/1i_hwcNFGmx_v72G_TZ9YjHjzUM6Yv74tvBlvb_CoHfU/edit#heading=h.bcsw102vr267) section of the design document.
