.PHONY: run build update install-autostart

run:
	./dist/twitch-redeem-trigger

build:
	mkdir -p dist
	rm -f dist/twitch-redeem-trigger
	go build -o dist/twitch-redeem-trigger src/main.go

update:
	go mod tidy
	go mod download
	go mod vendor

install-autostart:
	./service/install.sh
	sudo systemctl start twitch-redeem
	sudo systemctl enable twitch-redeem

# local extensions
-include Makefile.local