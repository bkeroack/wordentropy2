wordentropy2
============

https://www.wordentropy.org

Pseudo-grammatical Passphrase Generator

To run:

```
$ go get
$ go build
$ sudo ./wordentropy2    #requires elevated privs to bind to port 443
```

Put TLS key and certificate in a subdirectory called "tls"--application is hard coded to look for "cert.key" and "cert-unified.pem", respectively.

To generate plots via Plot.ly, put your credentials in data/plotly_creds.txt (exactly two lines: username and api key).
