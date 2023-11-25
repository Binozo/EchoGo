#!/bin/bash

export hostPath=$(pwd)
echo "Exporting to $hostPath..."
docker build --tag echoltrace:latest . &&
docker run -v "$hostPath":/ltrace/build echoltrace:latest &&
echo "Finished"