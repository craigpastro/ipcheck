# IP Checker

## How to run

Please make sure that you have [Go installed](https://golang.org/).

Start up Postgres:
```
docker run --rm -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:13.3
```
Then
```
go run .
```
should start the service listening on 127.0.0.1:8080.

You can configure various setting by setting environment variables in `.env`.

## Configuring the blocklist

The ipsets that the service uses can be configured by setting the `IP_SETS` variable in `.env`. By
default it is using the lists found in
[this example](https://github.com/firehol/blocklist-ipsets#adding-the-ipsets-in-your-fireholconf)
(which is a subset of the
[Level 1](https://github.com/firehol/blocklist-ipsets#level-1---basic) and
[Level 2](https://github.com/firehol/blocklist-ipsets#level-2---essentials) lists).

N.B. The most ipsets you add, the more the performance may suffer.

## References

- [Design document](https://docs.google.com/document/d/1i_hwcNFGmx_v72G_TZ9YjHjzUM6Yv74tvBlvb_CoHfU/edit?usp=sharing)
