all: build

build:
	go build -o lunafan-control main.go

clean:
	rm -f lunafan-control *.pkg.tar.zst
	rm -rf pkg src

package:
	makepkg -f

install: build
	sudo install -Dm755 lunafan-control /usr/bin/lunafan-control
	sudo install -Dm644 config.json /etc/lunafan-control/config.json
	sudo mkdir -p /etc/lunafan-control/configs
	sudo cp config.json /etc/lunafan-control/configs/default.json
	sudo install -Dm644 lunafan-control.service /usr/lib/systemd/system/lunafan-control.service
	sudo systemctl daemon-reload

.PHONY: all build clean package install
