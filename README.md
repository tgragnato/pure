
# Proxy

This is a proxy for TLS and HTTP.
Ingress traffic is captured by nftables.

Most HTTP traffic is upgraded to HTTPS with a redirect.

The SNI field of the TLS is read.
The connection is forwarded to the local socks proxy, or handled directly if a bypass is set.

Filtering is performed at the DNS and SNI layers.

It also able to handle IMAPS and APNS traffic of Apple devices, WhatsApp and SMB.
It's open source, but addresses and strings are hardcoded, you'll need to adapt the codebase to your needs.

## nftables

```
chain prerouting {
  type nat hook prerouting priority 0; policy accept;
  [...]
  iifname <ifname> ip saddr <ranges> ip daddr != <iface_ip> tcp dport domain redirect to :domain
  iifname <ifname> ip saddr <ranges> ip daddr != <iface_ip> udp dport domain redirect to :domain
  iifname <ifname> ip saddr <ranges> ip daddr != <iface_ip> udp dport ntp redirect to :ntp
  iifname <ifname> ip saddr <ranges> ip daddr != <iface_ip> tcp dport http redirect to :http
  iifname <ifname> ip saddr <ranges> ip daddr != <iface_ip> tcp dport https redirect to :https
  iifname <ifname> ip saddr <ranges> ip daddr <iface_ip> tcp dport microsoft-ds redirect to :https
  iifname <ifnane> ip saddr <ranges> ip daddr @aws tcp dport 453 redirect to :https
  iifname <ifname> ip saddr <ranges> ip daddr @apple tcp dport imaps redirect to :https
  iifname <ifname> ip saddr <ranges> ip daddr @fb tcp dport 5222 redirect to :https
  iifname <ifname> ip saddr <ranges> ip daddr @apple tcp dport 5223 redirect to :https
  [...]
  <IPv6>
  [...]
}

chain output_<wan> {
  [...]
  tcp sport 1025-65535 tcp dport https skuid proxy accept
  [...]
}

```
