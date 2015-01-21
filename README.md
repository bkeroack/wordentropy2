wordentropy2
============

https://www.wordentropy.org

Pseudo-grammatical Passphrase Generator (using [libwordentropy](https://github.com/bkeroack/libwordentropy))

To run:

```
$ go get
$ go build
$ sudo ./wordentropy2    #requires elevated privs to bind to port 443
```

Put TLS key and certificate in a subdirectory called "tls"--application is hardcoded to look for "cert.key" and "cert-unified.pem", respectively.

To generate plots via Plot.ly, put credentials in data/plotly_creds.txt (exactly two lines: username and api key). Make sure you have Python (2.7), numpy and plotly installed.

To turn off plot generation, use "-local=true" when running. Missing plots will break the how-random page. This option will also cause the application to bind to port 4343, which does not require root privileges.
