.PHONY: release

build:
	go build -ldflags "-H=windowsgui" -o logovnik.exe ./main.go

release:
	rm -rf ./release
	rm -rf ./winres
	rm -f rsrc_windows_386.syso
	rm -f rsrc_windows_amd64.syso
	mkdir release
	mkdir release/frontend
	cp ./frontend/index.html release/frontend
	go install github.com/tc-hib/go-winres@latest
	go-winres init
	cp ./icon/icon.png ./winres
	cp ./icon/icon16.png ./winres
	go-winres simply --icon winres/icon.png
	go-winres make
	go build -ldflags "-H windowsgui" -o release/logovnik.exe