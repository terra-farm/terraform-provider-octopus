default: build test

fmt:
	go fmt github.com/DimensionDataResearch/terraform-octopus/...

build: fmt
	go build -o _bin/terraform-provider-octopus
	cp _bin/terraform-provider-octopus _bin/terraform-provisioner-octopus

test: fmt
	go test -v github.com/DimensionDataResearch/terraform-octopus/vendor/octopus
