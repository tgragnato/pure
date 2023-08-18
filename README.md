# Pure

`pure` is a collection of software and tools that I use to manage my network and my storage.

Directory tree:
- `shshd` is a http proxy tailor-made to support the "signed hash protocol"
    [https://www.theiphonewiki.com/wiki/SHSH](https://www.theiphonewiki.com/wiki/SHSH)
- `sntpd` is a sntp v4 server that relays the time of the system on which it is running
    [https://www.rfc-editor.org/rfc/rfc2030](https://www.rfc-editor.org/rfc/rfc2030)
- `snid` is a proxy for TLS that does not terminate
    - it refuses non-TLS traffic and connections to bare IP addresses
    - checks the geographical location of the remote endpoints
    - logs the [Server Name Indication](https://www.rfc-editor.org/rfc/rfc3546.html)
- the rest of the project is just basic ansible