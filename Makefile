sha = $$(git rev-parse HEAD)
image = k8s-event-logger
builddir = /app/src/github.com/wfarr/$(image)

default:
	docker build -t $(image) -f Dockerfile.buildenv .
	docker run --rm -e GOPATH=/app -v $(PWD):$(builddir) -w $(builddir) $(image):latest go build -o ./bin/k8s-event-logger .

build-release: default
	docker build -t dubs/$(image):$(sha) -f Dockerfile.release .

release: build-release
	docker tag dubs/$(image):$(sha) dubs/$(image):latest
	docker push dubs/$(image):$(sha)
	docker push dubs/$(image):latest

deploy:
	SHA=$(sha) erb deployments/k8s-event-logger.yaml.erb | kubectl apply -f -
