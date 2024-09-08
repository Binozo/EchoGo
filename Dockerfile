FROM ghcr.io/binozo/echogosdkcompiler:latest

# Target file for compilation. Your main fun needs to be there
ENV TARGET="main.go"

LABEL org.opencontainers.image.source = "https://github.com/Binozo/EchoGoSDK"

# Compile our stuff
WORKDIR /EchoGoSDK

# For alsa but this doesn't work because of missing alsa config files
#RUN apt-get update
#RUN apt-get install libasound2-dev -y
#
#COPY binaries/libasound.so ${NDK_LIBS}
#RUN cp -r /usr/include/alsa ${NDK_INCLUDE}/

# Configured for Echo Dot 2. Gen
CMD env GOOS=android GOARCH=arm GOARM=7 CGO_ENABLED=1 go build -o build/main cmd/$TARGET