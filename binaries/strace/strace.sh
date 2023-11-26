#!/bin/bash

export hostPath=$(pwd)
echo "Exporting to $hostPath..."
docker build --tag echostrace:latest . &&
docker run -v "$hostPath":/strace/build echostrace:latest &&
echo "Finished"