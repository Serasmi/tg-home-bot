.PHONY: lint
lint:
	golangci-lint run --config=.golangci.yml

.PHONY: test
test:
	go test -v -race -count=1 ./...

.PHONY: deploy
deploy:
	./scripts/deploy.sh
