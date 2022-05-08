
HTTP/SNI proxy daemon for OpenBSD

```
#!/bin/ksh

daemon="/usr/local/bin/proxy"
daemon_user="nobody"

. /etc/rc.d/rc.subr

rc_bg=YES
rc_reload=NO

rc_cmd $1
```