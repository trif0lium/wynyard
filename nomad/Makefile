.PHONY: *

vm:
	apt update
	apt install software-properties-common wget gpg coreutils -y
	wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
	echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
	sudo apt-get update && sudo apt-get install nomad
	curl -fsSL get.docker.com | bash

wynyard-0: vm
	cp wynyard-0.nomad.service /etc/systemd/system/nomad.service
	systemctl enable nomad
	systemctl start nomad
	vgcreate vg0 /dev/sdb
	lvcreate --thin --size 45G --chunksize 256K --poolmetadatasize 1G vg0/lv0

wynyard-1: vm
	cp wynyard-1.nomad.service /etc/systemd/system/nomad.service
	systemctl enable nomad
	systemctl start nomad
	vgcreate vg0 /dev/sdb
	lvcreate --thin --size 45G --chunksize 256K --poolmetadatasize 1G vg0/lv0
