all: sntpd shshd snid dnsd nfguard
	ansible-playbook --ask-vault-password -i inventory.yaml playbook.yaml

clean:
	rm -f roles/firewall/files/sntpd
	rm -f roles/firewall/files/shshd
	rm -f roles/firewall/files/snid
	rm -f roles/firewall/files/dnsd
	rm -f roles/firewall/files/nfguard

sntpd: roles/firewall/files/sntpd

roles/firewall/files/sntpd:
	GOOS=linux GOARCH=amd64 go build -C sntpd -o ../roles/firewall/files/sntpd

shshd: roles/firewall/files/shshd

roles/firewall/files/shshd:
	GOOS=linux GOARCH=amd64 go build -C shshd -o ../roles/firewall/files/shshd

snid: roles/firewall/files/snid

roles/firewall/files/snid:
	GOOS=linux GOARCH=amd64 go build -C snid -o ../roles/firewall/files/snid

dnsd: roles/firewall/files/dnsd

roles/firewall/files/dnsd:
	GOOS=linux GOARCH=amd64 go build -C dnsd -o ../roles/firewall/files/dnsd

nfguard: roles/firewall/files/nfguard

roles/firewall/files/nfguard:
	GOOS=linux GOARCH=amd64 go build -C nfguard -o ../roles/firewall/files/nfguard