all: nfguard snowflake client magnetico grafana
	ansible-playbook --ask-vault-password -i inventory.yaml playbook.yaml

clean:
	rm -f roles/master/files/nfguard
	rm -f roles/master/files/snowflake
	rm -f roles/master/files/client
	rm -f roles/master/files/magnetico

pure_path := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

nfguard: roles/master/files/nfguard
	ansible-playbook -i inventory.yaml playbook.yaml --tags nfguard

roles/master/files/nfguard:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C cmd/nfguard -o $(pure_path)/roles/master/files/nfguard

snowflake: roles/master/files/snowflake
	ansible-playbook -i inventory.yaml playbook.yaml --tags snowflake

roles/master/files/snowflake:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C snowflake/proxy -o $(pure_path)/roles/master/files/snowflake

client: roles/master/files/client
	ansible-playbook -i inventory.yaml playbook.yaml --tags tor

roles/master/files/client:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C snowflake/client -o $(pure_path)/roles/master/files/client

magnetico: roles/master/files/magnetico
	ansible-playbook -i inventory.yaml playbook.yaml --tags magnetico

roles/master/files/magnetico:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C magnetico/. -o $(pure_path)/roles/master/files/magnetico

grafana: roles/master/files/grafana-12.0.2.linux-amd64.tar.gz
	ansible-playbook -i inventory.yaml playbook.yaml --tags grafana

roles/master/files/grafana-12.0.2.linux-amd64.tar.gz:
	curl https://dl.grafana.com/oss/release/grafana-12.0.2.linux-amd64.tar.gz -o roles/master/files/grafana-12.0.2.linux-amd64.tar.gz

pyroscope: roles/master/files/pyroscope_1.13.1_linux_amd64.tar.gz
	ansible-playbook -i inventory.yaml playbook.yaml --tags pyroscope

roles/master/files/pyroscope_1.13.1_linux_amd64.tar.gz:
	curl -L https://github.com/grafana/pyroscope/releases/download/v1.13.1/pyroscope_1.13.1_linux_amd64.tar.gz -o roles/master/files/pyroscope_1.13.1_linux_amd64.tar.gz
