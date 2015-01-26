wordentropy2
============

https://www.wordentropy.org

Pseudo-grammatical Passphrase Generator (using [libwordentropy](https://github.com/bkeroack/libwordentropy))

To run:

```
$ go get
$ go build
$ ./wordentropy2 -local
```

**Local mode**

- Does not generate plots/statistics on startup ("How Random?" page will have blank areas).
- Compiles views on every request.
- Binds to high, non-privileged port.

**TLS (SSL)**

Wordentropy requires the use of TLS (SSL). For testing purposes you can create a self-signed certificate and private key with the
following commands:

```bash
$ cd wordentropy2
$ mkdir ./tls && cd ./tls
$ openssl req -x509 -newkey rsa:4096 -keyout cert.key -out cert-unified.pem -days 365 -nodes
```

Application is hardcoded to look for "tls/cert.key" and "tls/cert-unified.pem" for the private key and certificate, respectively.

**Plots/Statistics**

To generate plots via Plot.ly, put credentials in data/plotly_creds.txt (exactly two lines: username and api key). Make sure you have Python (2.7), numpy and plotly installed. Run wordentropy2 in production mode.


