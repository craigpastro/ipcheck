# IP Checker

## TODO

[ ] Tests
[ ] Get `updateBlocklist` to run every 24.3 hours
[ ] Benchmark
[ ] Review error handling
[ ] If benchmarks are terrible find a better way to populate the table
[ ] Is the way I update the table sane?

## Usage

### Install Go

https://golang.org/

### Start Postgres

```
docker run --rm -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:13.3
```

### Configuration

Set any variables in [`.env`](./env) to configure the behaviour of the service.

- `SERVER_ADDR` - The address the server will listen on. The default is `127.0.0.1:8080`.
- `GIN_MODE` - The mode Gin will run in. Choose either `debug` or `release`. The default is `debug`.
- `DATABASE_URL` - The hostname address of a running Postgres instance. The default is `postgres://postgres:password@127.0.0.1:5432/postgres`. (This default works with the docker command above.)
- `ALL_MATCHES` - When an IP address appears on a blocklist, if this is set to `true` it will return all the blocklists that the IP address appears in. If `false` it will return the first instance found. Note that performance may be affected if set to `true`. The default is `false`.
- `IP_SETS_DIR` - The directory where you want to download the ipsets to. Default is `/tmp/ipsets`.
- `IP_SETS` - A comma separated list of blocklists that you would like the system to use. Please find the list of all blocklists at https://github.com/firehol/blocklist-ipsets. Note, adding many blocklists may lead to performance degradation. The default is `feodo.ipset,palevo.ipset,sslbl.ipset,zeus.ipset,zeus_badips,dshield.netset,spamhaus_drop.netset,spamhaus_edrop.netset,fullbogons.netset,openbl.ipset,blocklist_de.ipset`. (The Level 1 and some of Level 2 blocklists described [here](https://github.com/firehol/blocklist-ipsets#which-ones-to-use).)

### Running the service

```
go run .
```
should start the service listening at `SERVER_ADDR`.

### Tests

TODO

## References

- [Design document](https://docs.google.com/document/d/1i_hwcNFGmx_v72G_TZ9YjHjzUM6Yv74tvBlvb_CoHfU/edit?usp=sharing)
