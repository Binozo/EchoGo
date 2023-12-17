FROM ghcr.io/binozo/echogosdkcompiler:latest

# Compile our stuff
WORKDIR /EchoGoSDK

# Configured for Echo Dot 2. Gen
CMD env GOOS=android GOARCH=arm GOARM=7 CGO_ENABLED=1 go build -o build/main cmd/main.go