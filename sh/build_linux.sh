GOOS=linux GOARCH=amd64 go build -o start.linux.amd64 $(ls *.go | grep -v test | tr '\n' ' ')
mv start.linux.amd64 bin/linux/
