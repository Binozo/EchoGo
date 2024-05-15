FROM ghcr.io/binozo/echogosdkcompiler:latest

# Target file for compilation. Your main fun needs to be there
ENV TARGET="main.go"

LABEL org.opencontainers.image.source = "https://github.com/Binozo/EchoGoSDK"

# Compile our stuff
WORKDIR /EchoGoSDK

# Configured for Echo Dot 2. Gen
CMD env GOOS=android GOARCH=arm GOARM=7 CGO_ENABLED=1 go build -o build/main cmd/$TARGET