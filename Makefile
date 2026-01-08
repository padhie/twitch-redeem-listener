.PHONY: run build update install-autostart

run:
	./twitch-redeem-trigger

build:
	rm twitch-redeem-trigger || true
	go build -o twitch-redeem-trigger main.go

update:
	go mod tidy
	go mod download
	go mod vendor

install-autostart:
	./service/install.sh
	sudo systemctl start twitch-redeem
	sudo systemctl enable twitch-redeem
