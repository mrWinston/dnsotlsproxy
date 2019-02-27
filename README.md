# Dns-over-TLS Proxy

This Program implements a DNS-Over-TLS Proxy service. It accepts DNS Queries
via TCP and UDP and forwards them to a DNSoTLS Service of your choice.

It uses Golang 1.11 with go-modules.


## Using the Service

User the provided `docker-compose.yml` to build and run the service:

```
$> docker-compose up
```

Afterwards, you can point your application to the app for dns resolution (for
example with dig) :

```
$> dig @127.0.0.1 -p 8053 google.de +short
172.217.23.163
```


## Configuring the Service

All configuration is done via environment variables. The following values can
be configured:

Env-Key              | Default   | Description
---                  | ---       | ---
UDP_PORT             | `53`      | The port on which the app should listen for UDP-Requests
UDP_ADDRESS          | `0.0.0.0` | The Address on which the app should listen for UDP-Requests
TCP_PORT             | `53`      | The port on which the app should listen for TCP-Requests
TCP_ADDRESS          | `0.0.0.0` | The Address on which the app should listen for TCP-Requests
DNS_UPSTREAM_PORT    | `853`     | The Port of the upstream DNS-over-TLS Server
DNS_UPSTREAM_ADDRESS | `1.1.1.1` | The Address of the upstream DNS-over-TLS Server (Defaults to Cloudflare's DNS)
LOG_LEVEL            | `error`   | The log level, Possible Values: `info`, `error`

The `docker-compose.yml` configures the application to use port 8053 for both
tcp and udp, and sets the log-level to `info`


### Adding a custom certificate

If you're running your own DNS-over-TLS Service using a self-signed
certificate, you may want to inject your certificate into the running docker
container. By default, the container includes the system certificates from the
image `golang:1.11` which should be the same included in debian stretch.

To add your own certificates, mount them in the following directory in the
container: (docker run command for illustration):

```
$> docker run -v /my/certificate/path.cert:/etc/ssl/certs/my_cert.cert
```

## Security and Usability Considerations

The Application doesn't provide dns-over-tls itself. That means, that whenever
this app is deployed, it should be kept as closely (network wise) to its
clients as possible. An attacker in the same network will still be able to
sniff the dns requests coming into the service.

Even though the Cloudflare Dns server does support TLS1.3, this service does
not, so it's always using TLS 1.2 for connections. It is also refusing
connections with TLS 1.1 and 1.0 to prevent TLS-Downgrade attacks.

## Improvements

As it stands right now, the app does not support caching in any way. If dns
entries should be chached it is up to the client to do so. Future version could
implement this feature, as it would improve performance considerably. 

Also, TLS Connections are not reused across requests. This leads to the app
performing the tls handshake every time a request hits. This could be improved
by reusing a tls connection based to the dns upstream based on the ID of the
dns request. This would also improve performance a bit. 

Code-wise, the app is currently duplicating quite a lot of code. The TCP and
UDP listeners share quite a lot of common structure. In future iterations, it
would be smart to outsource some common logic into a separate component. 

Also, the server implementation leaves room for improvement. The UDP listener
does not support concurrent requests at the moment, while the TCP
implementation does to some extend. I would  still not use this in a high load
scenario.

In case the Upstream DNS Server uses padding in their responses, said padding
is NOT removed from the dns responses. Right now, i don't see any downside in
doing so. If, in the future this shows to be a problem, the padding would have
to be stripped and the Padding metadata (see
[here](https://tools.ietf.org/html/rfc7830)) would have to be removed. 
