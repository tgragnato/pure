# Pure

`pure` is a collection of software and tools that I use to manage my network and my storage.

Directory tree:
- `cmd/dnsd` is meant to be an intercepting DNS caching proxy
    - I'm willingly not implementing support for most record types
    - The forwarder for A and AAAA records is loosely inspired by
        - [Introducing DNS Resolver for Tor](https://blog.cloudflare.com/welcome-hidden-resolver/)
        - [DoHoT: better security, privacy, and integrity via load-balanced DNS over HTTPS over Tor](https://blog.apnic.net/2021/09/28/dohot-better-security-privacy-and-integrity-via-load-balanced-dns-over-https-over-tor/)
    - The cache is in-memory and it never performs record eviction, the purpose is to
        - minimize the latency introduced by DoHoT
        - reduce the amount of information deducible from the logs via pattern of life analysis
- `pkg/nfqueue` is a network interference tool that uses NFQUEUE
    - It modifies the window size of TCP packets with the SYN or the ACK flag set
    - Middleboxes that do not support stream reassembly will be unable to collect the metadata of the TLS streams
    - References:
        - [brdgrd (Bridge Guard)](https://github.com/NullHypothesis/brdgrd)
        - [About Geneva - How it Works](https://geneva.cs.umd.edu/about/)
- `pkg/sni` is a proxy for TLS that does not terminate
    - it refuses non-TLS traffic and connections to bare IP addresses
    - checks the geographical location of the remote endpoints
    - logs the [Server Name Indication](https://www.rfc-editor.org/rfc/rfc3546.html)
- `pkg/shsh` is a http proxy tailor-made to support the "signed hash protocol"
    - [https://www.theiphonewiki.com/wiki/SHSH](https://www.theiphonewiki.com/wiki/SHSH)
    - it also forwards the traffic of some [OCSP](https://www.rfc-editor.org/rfc/rfc6960) hosts
    - and redirects everything else towards HTTPS
- `pkg/sntp` is a sntp v4 server that relays the time of the system on which it is running
    - [https://www.rfc-editor.org/rfc/rfc2030](https://www.rfc-editor.org/rfc/rfc2030)
- the rest of the project is just basic ansible