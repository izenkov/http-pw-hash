# http-pw-hash

Hash Password String

**http-pw-hash** is a simple HTTP API server written in [Google Go](https://go.dev) as a learning exercise on how to build a production quality HTTP API server using only standard Go libraries.

## Repository restrictions

Internal hash repository is a simple slice of structs without any persistency. If you want to use this code as a starting point for your own server you would most likely need to adapt HashRepository to suit your needs. In real life you would most likely define a HashRepository interface, so you can plug in different implementations. I did not do it here for the sake of keeping the code as simple as possible.

## Endpoints

### POST /hash

Hash password

```sh
curl --data 'password=madCow' http://localhost:8080/hash
```

Call returns hash id as an integer. You need to store the id to retrieve actual hash in the following **/hash/{id}** call.

> Internal Id generator increments internal integer
> counter on every **/hash** call. If you call it with
> the same password you will get a unique integer id
> back on each call.
>
> In real application it is almost always better to use
> UUID as a unique id to avoid synchronization issues between
> threads, processes, and different servers. Also
> integer ids are very inconvenient if you need to merge
> ids from different servers, data centers, or regions.

### GET /hash/{id}

Get password hash by hash ID

```sh
curl http://localhost:8080/hash/17
```

There is a 5 seconds delay before you can get the actual hash value for your password. It simulates background processing (DB call, 3rd party services, etc.). If you make a call before the delay expires, you will get HTTP 202 back and ETA in milliseconds in the response message body. If you make a call after the delay you will get back SHA512 hash value of your password as Base64 string as your response message body.

> Calculating ETA it not always possible, so returning
> something like 'try it again later' is an alternative.
> ETA calculation in this project is trivial.

### GET /stats

Get password hash statistics

```sh
curl http://localhost:8080/stats
```

Returns JSON object with two fields. Total field returns a total number of **/hash** requests. Average field returns an average time for **/hash** call in nanoseconds.

> The average performance counter is recalculated on
> every **/hash** call. Sometimes it is better to recalculate
> performance counters on timer asynchronously.
> It depends on how expensive performance counters are
> from CPU utilization point of view. Only benchmarking
> can tell...

### POST | GET /shutdown

Gracefully shutdown HTTP server.

```sh
curl http://localhost:8080/shutdown
```

Very important feature in a production environment. It makes servers hot swap pain free.

> I know a few cases where a graceful server
> shutdown was not implemented in a production
> environment. Or it was implemented, but never tested
> under heavy load. This server shuts down gracefully
> under any load and it is [testable](/tests/README.md).

## Command line arguments

There is only one argument supported and it is the port to listen to.

```sh
./http-pw-hash 8787
```

The port argument is optional and if you don't specify it, it defaults to 8080.

> Tests scripts always assume server running on local
> host, default port 8080.

## Error handling

User input is validated and in case of bad input you will get HTTP 400 back and error details as string in the response message body.

## Tests

See [README.md](/tests/README.md) in [tests](/tests/) subfolder.

## Go version

Server was written using [Go 1.17.6](https://go.dev/dl)

## License

[MIT](/LICENSE)
