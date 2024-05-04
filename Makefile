all: dnsd nfguard snowflake client magnetico grafana
	ansible-playbook --ask-vault-password -i inventory.yaml playbook.yaml

clean:
	rm -f roles/firewall/files/dnsd
	rm -f roles/firewall/files/nfguard
	rm -f roles/firewall/files/snowflake
	rm -f roles/firewall/files/client
	rm -f roles/services/files/magneticod
	rm -f roles/services/files/magneticow

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

magnetico: magneticod magneticow

magneticod: roles/services/files/magneticod

roles/services/files/magneticod:
	GOOS=linux GOARCH=amd64 go build -C magnetico/cmd/magneticod --tags fts5 -o ../../../roles/services/files/magneticod

magneticow: roles/services/files/magneticow

roles/services/files/magneticow:
	GOOS=linux GOARCH=amd64 go build -C magnetico/cmd/magneticow --tags fts5 -o ../../../roles/services/files/magneticow

grafana: roles/services/files/grafana-10.4.2.linux-amd64.tar.gz

roles/services/files/grafana-10.4.2.linux-amd64.tar.gz:
	curl https://dl.grafana.com/oss/release/grafana-10.4.2.linux-amd64.tar.gz -o roles/services/files/grafana-10.4.2.linux-amd64.tar.gz
