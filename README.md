# Pure

`pure` is a collection of software and tools that I use to manage my network and my storage.

Directory tree:
- `pkg/dnshandlers` is meant to be an intercepting DNS caching proxy
    - I'm willingly not implementing support for most record types
- `pkg/dohot` is a forwarder for mainly A and AAAA records is loosely inspired by
    - [Introducing DNS Resolver for Tor](https://blog.cloudflare.com/welcome-hidden-resolver/)
    - [DoHoT: better security, privacy, and integrity via load-balanced DNS over HTTPS over Tor](https://blog.apnic.net/2021/09/28/dohot-better-security-privacy-and-integrity-via-load-balanced-dns-over-https-over-tor/)
- `pkg/shsh` is a http proxy tailor-made to support the "signed hash protocol"
    - [https://www.theiphonewiki.com/wiki/SHSH](https://www.theiphonewiki.com/wiki/SHSH)
    - it redirects everything else towards HTTPS
- `pkg/sntp` is a sntp v4 server that relays the time of the system on which it is running
    - [https://www.rfc-editor.org/rfc/rfc2030](https://www.rfc-editor.org/rfc/rfc2030)
- the rest of the project is just basic ansible