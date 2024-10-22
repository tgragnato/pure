all: dnsd nfguard snowflake client magnetico grafana spamd
	ansible-playbook --ask-vault-password -i inventory.yaml playbook.yaml

clean:
	rm -f roles/firewall/files/nfguard
	rm -f roles/firewall/files/snowflake
	rm -f roles/firewall/files/client
	rm -f roles/firewall/files/spamd
	rm -f roles/services/files/magnetico

nfguard: roles/firewall/files/nfguard

roles/firewall/files/nfguard:
	GOOS=linux GOARCH=amd64 go build -C cmd/nfguard -o ../../roles/firewall/files/nfguard

snowflake: roles/firewall/files/snowflake

roles/firewall/files/snowflake:
	GOOS=linux GOARCH=amd64 go build -C snowflake/proxy -o ../../roles/firewall/files/snowflake

client: roles/firewall/files/client

roles/firewall/files/client:
	GOOS=linux GOARCH=amd64 go build -C snowflake/client -o ../../roles/firewall/files/client

spamd: roles/firewall/files/spamd

roles/firewall/files/spamd:
	GOOS=linux GOARCH=amd64 go build -C cmd/spamd -o ../../roles/firewall/files/spamd

magnetico: roles/services/files/magnetico

roles/services/files/magnetico:
	GOOS=linux GOARCH=amd64 go build -C magnetico/. -o ../roles/services/files/magnetico

grafana: roles/services/files/grafana-11.2.3.linux-amd64.tar.gz

roles/services/files/grafana-11.2.3.linux-amd64.tar.gz:
	curl https://dl.grafana.com/oss/release/grafana-11.2.3.linux-amd64.tar.gz -o roles/services/files/grafana-11.2.3.linux-amd64.tar.gz
