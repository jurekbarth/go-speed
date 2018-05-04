# GoSpeed

Simple Fileserver with Speedmanipulation Options and HTTPS.

## Usage

### Define Ports

By default gospeed will listen on port 8080 and 8443, but you can change with setting flags. In order to change the http port use `-http=3000`, for https use `-https=3001`

For example `speed -http=3000 -https:3001`

### Delay Response

In order to delay response of a given resource use a query param e.g. `http://localhost:8080?speed=1`
The number passed to speed is multiplied by one second. So the example above will delay the response for one second.

### Generate your own Certs

```
# Generate key with algorithm "RSA" ≥ 2048-bit
openssl genrsa -out server.key 2048

# Generate public key based on private key.
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```
