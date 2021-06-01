export CGO_ENABLED=0
export DOCKER_BUILDKIT=1

.PHONY: kind
kind:
	-kind create cluster --name=cmstore

.PHONY: image
image:
	docker build -t cmstore .

.PHONY: build
build:
	go build -o bin/cmstore .

.PHONY: run
run: kind image
	kind load docker-image --name=cmstore cmstore
	kubectl --context=kind-cmstore apply -f _examples/manifests/rbac.yaml
	kubectl --context=kind-cmstore create -f _examples/manifests/pod.yaml

.PHONY: clean
clean:
	-kind delete cluster --name=cmstore