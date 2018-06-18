#!/bin/sh
set -e

# Designed for Windows Subsystem for Linux, though probably adaptable to most things.
# MAKE SURE THIS DRIVE LETTER IS CORRECT.
# RUN THIS ONLY ONCE.

LETTER=e

sudo mkdir /mnt/$LETTER
sudo mount -t drvfs $LETTER: /mnt/$LETTER

touch /mnt/$LETTER/ssh
echo "dtoverlay=dwc2" >> /mnt/$LETTER/config.txt
sudo sed -i 's/rootwait/rootwait modules-load=dwc2,g_ether/' /mnt/$LETTER/cmdline.txt
cp wpa_supplicant.conf /mnt/$LETTER

# CORRECTLY DISMOUNTING IT IS HIGHLY ENCOURAGED.
# Windows will not let you eject the card until you do so.
sudo umount /mnt/$LETTER
