# goreplay-http-logger

I needed a way to directly capture http traffic for use with [GoReplay](https://goreplay.org/), and there did not seem to be an official method. As such, I wrote my own quick server to do the job. I may as well share it with the world as it has been useful to me. I did not do anything fancy here, just a simple cli argument configuration.

```
$ ./goreplay-http-logger --help
http log server
  -bind string
        HTTP bind address
  -log-file string
        Log file name with date (default "http-%Y%m%d.log")
  -port int
        HTTP port (default 8080)
```
