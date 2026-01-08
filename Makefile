.PHONY: run build update

run:
	./twitch-redeem-trigger

build:
	rm twitch-redeem-trigger || true
	go build -o twitch-redeem-trigger main.go

update:
	go mod tidy
	go mod download
	go mod vendor
