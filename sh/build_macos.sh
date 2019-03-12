go build -o start $(ls *.go | grep -v test | tr '\n' ' ')
mv start bin/macos/
