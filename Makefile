.PHONY: install test build serve clean pack deploy ship undeploy

TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)

export TAG

install:
	go get .

test:
	go test ./...

build: install
	CGO_ENABLED=0 go build -ldflags "-X main.version=$(TAG)" -o news .

serve: build
	./news

clean:
	rm ./news

pack:
	GOOS=linux make build
	docker build -t pmichali/news-service:$(TAG) .

upload:
	docker push pmichali/news-service:$(TAG)

deploy:
	envsubst < k8s/deployment.yml | kubectl apply -f -

ship: test pack upload deploy

undeploy:
	kubectl delete -f k8s/deployment.yml
