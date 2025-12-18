release:
	GOOS=linux GOARCH=amd64; go build -ldflags="-s -w" -o scrapper
	tar -cJf scrapper.tar.xz scrapper
	rm -f scrapper
