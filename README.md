
# Proxy

Ingress -> nftables -> proxy -> Services

Egress <- proxy <- nftables <- LAN/VPN


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
