
# HAProxy as egress controller

Inspired by:
- https://www.haproxy.com/user-spotlight-series/haproxy-as-egress-controller/
- https://www.slideshare.net/roidelapluie/haproxy-as-egress-controller


## Traffic flows

Ingress -> HAProxy -> Services

Egress <- httpproxy <- HAProxy <- nftables <- LAN/VPN

Egress <- sniproxy <- HAProxy <- nftables <- LAN/VPN

## HAProxy

```
frontend http-in
  bind *:80
  use_backend http-egress if { src <ranges> }
  [...]

backend http-egress
  tcp-request content silent-drop unless { src <ranges> }
  tcp-request content reject unless HTTP
  server httpproxy 127.0.0.1:9080

listen tls
  bind *:443
  mode tcp

  stick-table type binary len 32 size 30k expire 30m

  acl clienthello req.ssl_hello_type 1
  acl serverhello res.ssl_hello_type 2

  tcp-request inspect-delay 10s
  tcp-request content reject unless clienthello

  stick on payload_lv(43,1) if clienthello
  stick store-response payload_lv(43,1) if serverhello

  [...]

  server sniproxy 127.0.0.1:9081 weight 0
  use-server sniproxy if { src <ranges> }

  [...]
```

## nftables

```
chain prerouting {
  type nat hook prerouting priority 0; policy accept;
  [...]
  iifname <ifname> ip saddr <ranges> ip daddr != <iface_ip> udp dport domain redirect to :domain
  iifname <ifname> ip saddr <ranges> ip daddr != <iface_ip> tcp dport http redirect to :http
  iifname <ifname> ip saddr <ranges> ip daddr != <iface_ip> tcp dport https redirect to :https
  [...]
}

chain output_<wan> {
  [...]
  tcp sport 1025-65535 tcp dport { http, https } skuid proxy accept
  [...]
}

```
