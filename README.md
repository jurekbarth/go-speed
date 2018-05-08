# GoSpeed

Simple Fileserver with Speedmanipulation Options and HTTPS.

## Usage

### Options

By default gospeed will listen on port 8080 and 8443 and uses a self signed certificate. You can pass a number of flags to customize it.

```
// options:
http: int // defines http port
https: int // defines https port
key: string // relative or absolute path to key file
cert: string // relative or absolute path to cert file
root: string // relative or absolute path to the servers root dir
defaultSpeed: float // default delays to server response times
```

For example `speed -http=3000 -https=3001 -key=server.key -cert=server.crt -root=./dist -defaultSpeed=0.25`

### Custom Delay Response

In order to delay response of a given resource use a query param e.g. `http://localhost:8080?speed=1`
The number passed to speed is multiplied by one second. So the example above will delay the response for one second.

### Generate your own Certs

```
# Generate key with algorithm "RSA" ≥ 2048-bit
openssl genrsa -out server.key 2048

# Generate public key based on private key.
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```
