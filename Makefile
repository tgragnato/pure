all: nfguard snowflake client magnetico grafana
	ansible-playbook --ask-vault-password -i inventory.yaml playbook.yaml

clean:
	rm -f roles/master/files/nfguard
	rm -f roles/master/files/snowflake
	rm -f roles/master/files/client
	rm -f roles/master/files/magnetico

nfguard: roles/master/files/nfguard

roles/master/files/nfguard:
	GOOS=linux GOARCH=amd64 go build -C cmd/nfguard -o ../../roles/master/files/nfguard

snowflake: roles/master/files/snowflake

roles/master/files/snowflake:
	GOOS=linux GOARCH=amd64 go build -C snowflake/proxy -o ../../roles/master/files/snowflake

client: roles/master/files/client

roles/master/files/client:
	GOOS=linux GOARCH=amd64 go build -C snowflake/client -o ../../roles/master/files/client

magnetico: roles/master/files/magnetico

roles/master/files/magnetico:
	GOOS=linux GOARCH=amd64 go build -C magnetico/. -o ../roles/master/files/magnetico

grafana: roles/master/files/grafana-11.5.0.linux-amd64.tar.gz

roles/master/files/grafana-11.5.0.linux-amd64.tar.gz:
	curl https://dl.grafana.com/oss/release/grafana-11.5.0.linux-amd64.tar.gz -o roles/master/files/grafana-11.5.0.linux-amd64.tar.gz

pyroscope: roles/master/files/pyroscope_1.12.0_linux_amd64.tar.gz

roles/master/files/pyroscope_1.12.0_linux_amd64.tar.gz:
	curl -L https://github.com/grafana/pyroscope/releases/download/v1.12.0/pyroscope_1.12.0_linux_amd64.tar.gz -o roles/master/files/pyroscope_1.12.0_linux_amd64.tar.gz
