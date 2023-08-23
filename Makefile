all: sntpd shshd snid
	ansible-playbook --ask-vault-password -i inventory.yaml playbook.yaml

clean:
	rm -f roles/firewall/files/sntpd
	rm -f roles/firewall/files/shshd
	rm -f roles/firewall/files/snid

sntpd: roles/firewall/files/sntpd

roles/firewall/files/sntpd:
	GOOS=linux GOARCH=amd64 go build -C sntpd -o ../roles/firewall/files/sntpd

shshd: roles/firewall/files/shshd

roles/firewall/files/shshd:
	GOOS=linux GOARCH=amd64 go build -C shshd -o ../roles/firewall/files/shshd

snid: roles/firewall/files/snid

roles/firewall/files/snid:
	GOOS=linux GOARCH=amd64 go build -C snid -o ../roles/firewall/files/snid