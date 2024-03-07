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

Example Nginx config for mirroring requests:

```nginx
upstream backend {
    server 127.0.0.1:8087;
}

upstream mirror_backend {
    server 127.0.0.1:8080;
}
server {
    listen       8086 default_server;
    server_name  localhost;

    # Send body to mirror.
    mirror_request_body on;

    # redirect server error pages to the static page /50x.html
    #
    error_page   500 502 503 504  /50x.html;
    location = /50x.html {
        root   /usr/share/nginx/html;
    }

    location / {
        mirror @mirror;
        proxy_pass http://backend;
    }

    location = @mirror {
        internal;
        proxy_pass http://mirror_backend$request_uri;
    }
}
```
