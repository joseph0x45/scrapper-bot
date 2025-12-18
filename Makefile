release:
	GOOS=linux GOARCH=amd64; go build -ldflags="-s -w"
	tar -cJf scrapper.tar.xz simon
	rm -f scrapper
