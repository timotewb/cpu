
sudo apt-get install gcc-arm*
GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc go build -o build/rss

GOOS=linux GOARCH=arm GOARM=7 CC=arm-linux-gnueabi-gcc go build -o build/rss

If building on machine running just do `go build -o`