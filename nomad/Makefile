.PHONY: *

wynyard-0:
	cp wynyard-0.nomad.service /etc/systemd/system/nomad.service
	systemctl enable nomad
	systemctl start nomad
	vgcreate vg0 /dev/sdb
	lvcreate --thin --size 45G --chunksize 256K --poolmetadatasize 1G vg0/lv0

wynyard-1:
	cp wynyard-1.nomad.service /etc/systemd/system/nomad.service
	systemctl enable nomad
	systemctl start nomad
	vgcreate vg0 /dev/sdb
	lvcreate --thin --size 45G --chunksize 256K --poolmetadatasize 1G vg0/lv0
