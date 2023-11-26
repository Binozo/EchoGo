#!/bin/bash

echo "Building compiler..."
docker build -f compiler/Dockerfile -t echogosdkcompiler:latest .

echo "Building project"
docker build --tag echogo:latest .
docker run -v "$(pwd)":/EchoGoSDK/build echogo:latest

echo "Building ltrace"
root=$(pwd)
cd binaries/ltrace/
bash ltrace.sh
cd $root

echo "Building strace"
cd binaries/strace/
bash strace.sh
cd $root