#!/bin/bash

if [ "$EUID" -ne 0 ]
  then echo "Please run this script as root"
  exit
fi

targetWorkingDir="/opt/echogo"
workingDir=$PWD
preloaderFileName="preloader_no_hdr.bin"
executableFileName="echogo"

scriptName=$(basename $BASH_SOURCE)

install () {
  # Check if preloader bin file exists
  if [ ! -f $preloaderFileName ]; then
      echo "$preloaderFileName not found! Please copy it into this directory!"
      exit
  fi
  # Check if your compiled program exists
  if [ ! -f $executableFileName ]; then
      echo "Your executable '$executableFileName' not found! Please copy it into this directory!"
      exit
  fi

  echo "Copying script..."
  rm -rf $targetWorkingDir
  mkdir -p $targetWorkingDir
  cp "$scriptName" $targetWorkingDir/
  chmod +x $targetWorkingDir/$scriptName # just in case
  echo "Copying your executable..."
  cp "$executableFileName" $targetWorkingDir/

  echo "Moving preloader file"
  mv $preloaderFileName "$targetWorkingDir/$preloaderFileName"

  # Install systemd service
  echo "Installing service..."
  rm /etc/systemd/system/echogo.service
  cat <<EOT >> /etc/systemd/system/echogo.service
[Unit]
Description=Boot system for the Echo Dot

[Service]
Type=simple
WorkingDirectory=/opt/echogo
ExecStart=/bin/bash /opt/echogo/boot.sh

[Install]
WantedBy=multi-user.target
EOT
  systemctl daemon-reload
  systemctl enable echogo

  # Build mtkclient
  echo "Building mtkclient..."
  apt install -y python3 git libusb-1.0-0 python3-pip adb
  cd "$targetWorkingDir/"
  git clone https://github.com/bkerler/mtkclient.git
  cd mtkclient
  pip3 install -r requirements.txt
  pip3 install .
  # Install rules
  usermod -a -G plugdev $USER
  usermod -a -G dialout $USER
  cp mtkclient/Setup/Linux/*.rules /etc/udev/rules.d
  udevadm control -R

  echo "Starting Service"
  systemctl start echogo
  echo "Started Service. Now you will get live logs. You can always ^C out"
  sleep 1 # give some time to read
  journalctl -u echogo --no-pager -f
}

askInstall () {
  while true; do
        read -p "This script will move itself to $targetWorkingDir for installation. Proceed? (y/n) " yn
        case $yn in
            [Yy]* ) install; break;;
            [Nn]* ) exit;;
            * ) echo "Please answer yes or no.";;
        esac
    done
}

if [ "$workingDir" != "$targetWorkingDir" ]; then
  askInstall
  exit
fi

echo "Booting Echo Dot"
python3 mtkclient/mtk plstage --preloader=preloader_no_hdr.bin
echo "Waiting for Device to come online..."
adb wait-for-device
echo "Uploading your executable..."
adb push $targetWorkingDir/$executableFileName /data/local/tmp/echogo
echo "Starting your executable..."
adb shell "chmod +x /data/local/tmp/echogo"

while true
do
  adb shell "cd /data/local/tmp/ && ./echogo"
  echo "The process died. Restarting..."
done