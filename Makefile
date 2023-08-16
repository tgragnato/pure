clean:
	rm roles/firewall/files/sntpd

sntpd: roles/firewall/files/sntpd

roles/firewall/files/sntpd:
	GOOS=linux GOARCH=amd64 go build -C sntpd -o ../roles/firewall/files/sntpd