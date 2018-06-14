#!/bin/sh
passwd
sudo apt update
sudo apt dist-upgrade
sudo raspi-config
echo "Suggestion: `sudo reboot` now."
