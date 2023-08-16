clean:
	rm roles/firewall/files/sntpd
	rm roles/firewall/files/shshd

sntpd: roles/firewall/files/sntpd

roles/firewall/files/sntpd:
	GOOS=linux GOARCH=amd64 go build -C sntpd -o ../roles/firewall/files/sntpd

shshd: roles/firewall/files/shshd

roles/firewall/files/shshd:
	GOOS=linux GOARCH=amd64 go build -C shshd -o ../roles/firewall/files/shshd