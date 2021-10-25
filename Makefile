IMAGE_REPO?=public.ecr.aws/micahhausler/route53-dynamic-dns
VERSION=$(shell cat VERSION)
IMAGE?=${IMAGE_REPO}:${VERSION}


bins:
	mkdir -p _output/bin/linux-arm64
	mkdir -p _output/bin/linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
		go build -trimpath -v -ldflags='-s -w --buildid=""' \
		-o _output/bin/linux-arm64/route53-dynamic-dns
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -trimpath -v -ldflags='-s -w --buildid=""' \
		-o _output/bin/linux-amd64/route53-dynamic-dns

# Set to "true" to build/push multi-arch
PUBLISH?=

IMAGE_OUTPUT=--load
IMAGE_PLATFORM= linux/amd64
ifeq ($(PUBLISH),true)
	IMAGE_OUTPUT= --push
	IMAGE_PLATFORM= linux/arm64,linux/amd64
endif
.PHONY: image
image: bins
	docker buildx build \
		--platform=$(IMAGE_PLATFORM) \
		$(IMAGE_OUTPUT) \
		--tag ${IMAGE} \
		.
