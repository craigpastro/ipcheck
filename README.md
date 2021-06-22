# IP Checker

## TODO

- [ ] Review error handling
- [ ] If benchmarks are terrible find a better way to populate the table
- [ ] Is the way I update the table sane?

## Design document

There is a [design document](https://docs.google.com/document/d/1i_hwcNFGmx_v72G_TZ9YjHjzUM6Yv74tvBlvb_CoHfU/edit#) which describes the general system, assumptions, tradeoffs, future work, etc.

## Usage

### Install Go

https://golang.org/

### Start Postgres

```
docker run --rm -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:13.3
```

### Configuration

Set any variables in [`.env`](./.env) to configure the behaviour of the service.

- `SERVER_ADDR` - The address the server will listen on. The default is `127.0.0.1:8080`.
- `GIN_MODE` - The mode Gin will run in. Choose either `debug` or `release`. The default is `release`.
- `DATABASE_URL` - The hostname address of a running Postgres instance. The default is `postgres://postgres:password@127.0.0.1:5432/postgres`. (This default works with the docker command above.)
- `ALL_MATCHES` - When an IP address appears on a blocklist, if this is set to `true` it will return all the blocklists that the IP address appears in. If `false` it will return the first instance found. Note that performance may be affected if set to `true`. The default is `false`.
- `IP_SETS_DIR` - The directory where you want to download the ipsets to. Default is `/tmp/ipsets` (which will not work on Windows so please update it if you are using Windows).
- `IP_SETS` - A comma separated list of blocklists that you would like the system to use. Please find the list of all blocklists at https://github.com/firehol/blocklist-ipsets. Note, adding many blocklists may lead to performance degradation. The default is `feodo.ipset,palevo.ipset,sslbl.ipset,zeus.ipset,zeus_badips,dshield.netset,spamhaus_drop.netset,spamhaus_edrop.netset,fullbogons.netset,openbl.ipset,blocklist_de.ipset`. (The Level 1 and some of Level 2 blocklists described [here](https://github.com/firehol/blocklist-ipsets#which-ones-to-use).)

### Running the service

```
go run .
```
should start the service listening at `SERVER_ADDR`.

### Tests

Start Postgres (note the port and name):
```
docker run --rm -d --name test_postgres -p 6543:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_DB=test postgres:13.3
```

Run the tests:
```
go test ./...
```

When the tests are finish remove the database with:
```
docker stop test_postgres
```

## Benchmarking

You can benchmark using a partial subset of the [blocklist_de.ipset](./test_ipsets/benchmark.ipset) (the first 10,000 entries of the 2021-06-19 version).

For example, using `siege` in benchmark mode with 10 concurrent users for 5 minutes:
```
$ siege -b -c10 -t5m -f test_ipsets/benchmark.ipset

{
	"transactions":			        110167,
	"availability":			        100.00,
	"elapsed_time":			        299.85,
	"data_transferred":		       13.19,
	"response_time":		          0.03,
	"transaction_rate":		      367.41,
	"throughput":			            0.04,
	"concurrency":			          9.96,
	"successful_transactions":  110167,
	"failed_transactions":		       0,
	"longest_transaction":		    1.12,
	"shortest_transaction":		    0.00
}
```

## Steps to make IP Checker production ready

Some steps are described in the [Future Work](https://docs.google.com/document/d/1i_hwcNFGmx_v72G_TZ9YjHjzUM6Yv74tvBlvb_CoHfU/edit#heading=h.bcsw102vr267) section of the design document.
