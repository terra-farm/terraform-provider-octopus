default: build

build:
	go build -o _bin/terraform-provider-octopus
	cp _bin/terraform-provider-octopus _bin/terraform-provisioner-octopus

test:
	go test -v github.com/DimensionDataResearch/terraform-octopus/vendor/octopus
