# Password Breach Checker

The purpose of this repository is a proof of concept to make a clone of the API found at: [haveibeenpwned](https://haveibeenpwned.com/Passwords)

The general guideline is to explore optimization techniques in Go and use the provided password hash file as much as possible to avoid using any database while at the same time having a high performant service.

To download the password hashes file with counts (ordered by hash) go to [haveibeenpwned](https://haveibeenpwned.com/Passwords) and follow the instructions. Alternatively, go directly into the [PwnedPasswordsDownloader repository](https://github.com/HaveIBeenPwned/PwnedPasswordsDownloader), install the tool and run it locally to get the latest version (around 35.55GB with 866'044,561 hash/count pairs at the time of writing).

Create a `data` folder at the root and save your downloaded `.txt` file there, the database package tests expect the name `data/pwned-passwords-sha1-ordered-by-hash-v8.bin`, but name it as you want if not running tests. To get the `.bin` file, compile and run `cmd\process` by passing the `.txt` file as an input with the `-f` flag. This will process the text file into binary with constant size per hash(20) + count(4) pair of 24 bytes. With the binary file, the database package can use **Memory Mapping** to perform **Binary Search** on it for a highly efficient, read-only, search engine.

The `cmd/check` utility can be used to quickly check a password against the database.

The `cmd/server` provides a quick API to check passwords against the database by receiving a SHA1 hash and returning a counter of breaches, where `0` means `Not Found` in the database. It includes a minimal frontend to test the API embedded in the binary. It has some flags to pass the binary database file, specify a port, or disable logging when running load tests.

## K6 Load Testing
A simple [K6](https://k6.io/) script is provided in the `load.js` file for Load Testing the API.

```plain
❯ k6 run --vus 30 --duration 60s load.js

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: load.js
     output: -

  scenarios: (100.00%) 1 scenario, 30 max VUs, 1m30s max duration (incl. graceful stop):
           * default: 30 looping VUs for 1m0s (gracefulStop: 30s)


running (1m00.0s), 00/30 VUs, 3012451 complete and 0 interrupted iterations
default ✓ [======================================] 30 VUs  1m0s

     data_received..................: 353 MB  5.9 MB/s
     data_sent......................: 603 MB  10 MB/s
     http_req_blocked...............: avg=1.41µs   min=0s med=0s      max=5.42ms  p(90)=0s     p(95)=0s
     http_req_connecting............: avg=12ns     min=0s med=0s      max=1.91ms  p(90)=0s     p(95)=0s
     http_req_duration..............: avg=450.23µs min=0s med=505.7µs max=26.5ms  p(90)=1ms    p(95)=1.09ms
       { expected_response:true }...: avg=450.23µs min=0s med=505.7µs max=26.5ms  p(90)=1ms    p(95)=1.09ms
     http_req_failed................: 0.00%   ✓ 0            ✗ 3012451
     http_req_receiving.............: avg=19.32µs  min=0s med=0s      max=24.52ms p(90)=0s     p(95)=0s
     http_req_sending...............: avg=10.09µs  min=0s med=0s      max=9.7ms   p(90)=0s     p(95)=0s
     http_req_tls_handshaking.......: avg=0s       min=0s med=0s      max=0s      p(90)=0s     p(95)=0s
     http_req_waiting...............: avg=420.81µs min=0s med=505µs   max=16.3ms  p(90)=1ms    p(95)=1.05ms
     http_reqs......................: 3012451 50206.965562/s
     iteration_duration.............: avg=592.1µs  min=0s med=539.5µs max=26.81ms p(90)=1.05ms p(95)=1.34ms
     iterations.....................: 3012451 50206.965562/s
     vus............................: 30      min=30         max=30
     vus_max........................: 30      min=30         max=30
```
