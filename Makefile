all: dnsd nfguard
	ansible-playbook --ask-vault-password -i inventory.yaml playbook.yaml

clean:
	rm -f roles/firewall/files/dnsd
	rm -f roles/firewall/files/nfguard

dnsd: roles/firewall/files/dnsd

roles/firewall/files/dnsd:
	GOOS=linux GOARCH=amd64 go build -C dnsd -o ../roles/firewall/files/dnsd

nfguard: roles/firewall/files/nfguard

roles/firewall/files/nfguard:
	GOOS=linux GOARCH=amd64 go build -C cmd/nfguard -o ../../roles/firewall/files/nfguard