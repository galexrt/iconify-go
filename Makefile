.PHONY: test

bench:
	go test -bench=. -benchmem

test:
	go test -coverprofile cover.out ./...

coverage: test
	go tool cover -html=cover.out -o cover.html
	open cover.html