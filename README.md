# Pure

`pure` is a collection of software and tools that I use to manage my network and my storage.

Directory tree:
- `shshd` is a http proxy tailor-made to support the "signed hash protocol"
    - [https://www.theiphonewiki.com/wiki/SHSH](https://www.theiphonewiki.com/wiki/SHSH)
    - it also forwards the traffic of some ocsp hosts
- `sntpd` is a sntp v4 server that relays the time of the system on which it is running
    - [https://www.rfc-editor.org/rfc/rfc2030](https://www.rfc-editor.org/rfc/rfc2030)
- `snid` is a proxy for TLS that does not terminate
    - it refuses non-TLS traffic and connections to bare IP addresses
    - checks the geographical location of the remote endpoints
    - logs the [Server Name Indication](https://www.rfc-editor.org/rfc/rfc3546.html)
- `dnsd` is meant to be an intercepting DNS caching proxy
    - I'm willingly not implementing support for most record types
    - The forwarder for A and AAAA records is loosely inspired by
        - [Introducing DNS Resolver for Tor](https://blog.cloudflare.com/welcome-hidden-resolver/)
        - [DoHoT: better security, privacy, and integrity via load-balanced DNS over HTTPS over Tor](https://blog.apnic.net/2021/09/28/dohot-better-security-privacy-and-integrity-via-load-balanced-dns-over-https-over-tor/)
    - The cache is in-memory and it never performs record eviction, the purpose is to
        - minimize the latency introduced by DoHoT
        - reduce the amount of information deducible from the logs via pattern of life analysis
- the rest of the project is just basic ansible