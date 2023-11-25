FROM echogosdkcompiler:latest

# Compile our stuff
WORKDIR /EchoGoSDK
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

# Configured for Echo Dot 2. Gen
CMD env GOOS=android GOARCH=arm GOARM=7 CGO_ENABLED=1 go build -o build/main example/AudioEcho.go