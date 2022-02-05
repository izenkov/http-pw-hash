# http-pw-hash Bash test scripts

> Tests assume that http-pw-hash server runs
> locally on port 8080

## Scripts

- **bulk-hash**
  - *passwords.txt*
- **loop-get**
- **loop-stats**
- **test-api**
  - test-api.conf

## bulk-hash

Stress test Bash script to hash multiple passwords by calling **/hash** endpoint.
Run with '-h' to see how to use it.

```sh
./bulk-hash -h
```

Script uses [passwords.txt](passwords.txt) containing 100000 unique most frequently used password.

## loop-get

Stress test Bash script to fetch password hashes by calling **/hash/{id}** endpoint. In one loop it gets all the existing hashes from first to last and from last to first.
Run with '-h' to see how to use it.

```sh
./loop-get -h
```

## loop-stats

Stress test Bash script to fetch server stats by calling **/stats** endpoint. In one loop it fetches one status.
Run with '-h' to see how to use it.

```sh
./loop-stats -h
```

## test-api

Bash script to test all endpoints. It has positive and negative tests. Run with '-h' to see how to use it.

```sh
./test-api -h
```

Script uses [test-api.conf](test-api.conf) containing various curl commands.

## Race stress test (go -race option)

For a best result run two instances of each script **bulk-hash**, **loop-get**, and **loop-stats** in parallel. WARNING: too many running scripts can render your machine unusable because of high CPU utilization.

## Graceful shutdown test

Launch instance of each script **bulk-hash**, **loop-get**, and **loop-stats** in parallel. Call **/shutdown** endpoint and observe script errors. If you see 'curl error 7' in all of your scripts instances it means that HTTP server shut down gracefully.

```sh
curl error: 7
```

If you see 'curl error 56' or similar curl error it means HTTP server did not shut down gracefully.

```sh
curl error: 56
```

> Obviously I don't expect you to see bad shutdown errors,
> but if you do, please let me know.

## Test environment

All test where performed on [Ubuntu 20.04 LTS](https://ubuntu.com)
