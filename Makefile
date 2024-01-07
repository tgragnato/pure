all: dnsd nfguard snowflake client
	ansible-playbook --ask-vault-password -i inventory.yaml playbook.yaml

clean:
	rm -f roles/firewall/files/dnsd
	rm -f roles/firewall/files/nfguard
	rm -f roles/firewall/files/snowflake
	rm -f roles/firewall/files/client

dnsd: roles/firewall/files/dnsd

roles/firewall/files/dnsd:
	GOOS=linux GOARCH=amd64 go build -C cmd/dnsd -o ../../roles/firewall/files/dnsd

nfguard: roles/firewall/files/nfguard

roles/firewall/files/nfguard:
	GOOS=linux GOARCH=amd64 go build -C cmd/nfguard -o ../../roles/firewall/files/nfguard

snowflake: roles/firewall/files/snowflake

roles/firewall/files/snowflake:
	GOOS=linux GOARCH=amd64 go build -C snowflake/proxy -o ../../roles/firewall/files/snowflake

client: roles/firewall/files/client

roles/firewall/files/client:
	GOOS=linux GOARCH=amd64 go build -C snowflake/client -o ../../roles/firewall/files/client
