# Raspberry Pi Guide

Clive Chan, June 2018

This assumes RPi 3B or RPi Zero W (mostly the latter), and Raspbian Stretch.
It might work on older RPis but I don't know. It probably doesn't work for previous versions of Raspbian.

## Flashing Raspbian and setting up headless
If you have problems with this, [this](https://gist.github.com/gbaman/975e2db164b3ca2b51ae11e45e8fd40a)
and [this](https://cdn-learn.adafruit.com/downloads/pdf/turning-your-raspberry-pi-zero-into-a-usb-gadget.pdf)
are great resources.

Flash the SD card with the latest version of Raspbian, following [these instructions](https://www.raspberrypi.org/documentation/installation/installing-images/).

Modify three files in the boot partition (you might need to remove and reinsert the SD card):

- Create `ssh`, a blank file.
- Add a new line to `config.txt` containing `dtoverlay=dwc2`.
- Add `modules-load=dwc2,g_ether` after `rootwait`. This file is very particular about spacing; delimit your insertion by a single space on each side.

Remember to use unix-style linebreaks where relevant.

Insert the SD card into the Pi, and plug in the Raspberry Pi to your laptop's USB port. *On RPiZero, make sure it's using the one labeled `USB`, not the one labeled `PWR IN`.*

Wait a bit for it to do a first-time boot-up. Meanwhile, if you're on Windows, install Bonjour Print Services.

Then to connect to it, run `ssh pi@raspberrypi.local`. [If you're on Windows, use PuTTY; it'll discover it for you and then you can use PuTTY local tunnels to ssh from whatever other terminal you want.]
[I haven't figured out how to make `avahi-daemon` discovery work on the Linux subsystem.]
It'll probably tell you that the authenticity of host 'raspberrypi' cannot be established, which is correct and inevitable; accept it.

(The RPi also helpfully provides its own hostname, `raspberrypi` to the AP which will put it into the DDNS if it's supported. Configuring over USB is a lot better than configuring over WiFi though, especially as some networks block p2p connections.)

The password is `raspberry`. Upon logging in, set the password with `passwd`. Please.

## Things to run on first boot

- As above, set the password with `passwd`.
- Connect to the internet somehow. See below for how to connect directly to any wifi network. [IN PROGRESS: usb tethering?]
- Run `sudo apt update` then `sudo apt dist-upgrade` on any new system to update all your packages. If you're wondering, [this](https://askubuntu.com/a/226213) is a good answer describing the difference between `upgrade` and `dist-upgrade`. It's probably a bad idea to do this on cell data, so see the Wifi section for how to do wifi.
- Run `sudo raspi-config`:
  - Activate the Camera interface if you're using that
  - The Expand Filesystem advanced option usually is good
  - Updating the raspi-config script is good too
- `sudo reboot` after all of this to clean up

## Debugging

`sudo journalctl`

## WiFi

Edit `/etc/wpa_supplicant/wpa_supplicant.conf`: (you can actually put wpa_supplicant.conf in the boot partition and it'll copy over on startup)

```
country=US
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1

network={
    ssid="WIFI_NAME"
    psk="WIFI_PASSWORD"
    key_mgmt=WPA-PSK
}
```

After modifying the file, run `sudo systemctl restart networking.service`
(there may be more specific commands; this one restarts the whole networking stack)
to reload the configuration changes. You might also just need a `sudo reboot` if it isn't working.

The above **doesn't work** on Orange Pi Zero. I instead appended this to `/etc/network/interfaces`:
```
auto wlan0
iface wlan0 inet dhcp
wpa-ssid WIFI_NAME
wpa-psk WIFI_PASSWORD
```


Useful commands:
- `iwconfig` for wifi specifically and `ifconfig` for network interfaces in general
- `iwgetid wlan0` will get you the currently connected wifi name.
- `sudo iw dev wlan0 scan | grep SSID` will scan for all nearby networks.
- `wpa_cli -i wlan0 reconfigure` will restart things after you edit `wpa_supplicant.conf`
- (not sure if this is necessary) `sudo mv /etc/ifplugd/action.d/action_wpa /etc/ifplugd/action.d/.action_wpa` if you're encountering issues where wlan0 gets cut off when eth gets connected
  - check this by running `ifconfig` and verifying that there's no ip address assigned to wlan0
- Other important files include `/etc/network/interfaces` and `/etc/dhcpcd.conf` but generally don't touch.

It's good practice to use only secured wifi (i.e. I like my phone hotspot), but if you really need to you can remove the `psk=...` line and use `key_mgmt=NONE` for WEP wifi. Don't ask about WPA-Enterprise (i.e. eduroam); it's really hard and eduroam also additionally refuses to let you talk to other computers on the same network, for obvious security reasons.

Alternatively, you can use `sudo raspi-config` for an even easier to use interface.

If you have lots of networks you can set `priority` for each one. By default it's `0`, and higher priority means it gets selected first. Negative is allowed.

## Protips

Use `sudo poweroff` before unplugging your Pi to gently shut down, avoiding bad things that usually don't happen but occasionally do.

By the way, you can (usually) move your SD card between different Raspberry Pi devices if you want.

## Setting up SSH

If this is going to be network-connected, please use SSH public key auth and not passwords.

Add your SSH key (probably in your computer's `~/.ssh/id_rsa.pub`) to its own line in `~/.ssh/authorized_keys`. Log out (`exit`) and try `ssh pi@raspberrypi` again. If it worked, it should not ask you for your password. If it didn't work, do not proceed until it works.

Now that ssh is working, edit (with `sudo`) the file `/etc/ssh/sshd_config`, and replace the line

```
#PasswordAuthentication yes
```

with

```
PasswordAuthentication no
```

for security. After restarting sshd, you will now be able to safely access your raspberry pi on whatever networks you want.

## Network security

For further security assurance, you'd probably like to know what's listening on your Pi. `sudo netstat -tuplna` will give you a list of all connections on the current device. By default it should output something like this:

```
Active Internet connections (servers and established)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN      516/sshd
tcp        0    304 192.168.43.139:22       192.168.43.235:51211    ESTABLISHED 699/sshd: pi [priv]
tcp6       0      0 :::22                   :::*                    LISTEN      516/sshd
udp        0      0 0.0.0.0:5353            0.0.0.0:*                           205/avahi-daemon: r
udp        0      0 0.0.0.0:68              0.0.0.0:*                           443/dhcpcd
udp        0      0 0.0.0.0:34011           0.0.0.0:*                           205/avahi-daemon: r
udp6       0      0 :::5353                 :::*                                205/avahi-daemon: r
udp6       0      0 :::53708                :::*                                205/avahi-daemon: r
```

Respectively, these are:

```
1) The SSH daemon that is listening for connections.
2) My current SSH session with the pi.
3) The SSH daemon also listens on IPv6.
4) Honestly, no idea what avahi-daemon is but I think it allows Bonjour to discover your device and locally resolve the DNS, instead of depending on the LAN's DNS. Most phone hotspots do DNS though. It's not actively listening so I don't really care that much. It's built into Raspbian.
5) DHCP, which gets you an IP address from the local network and a bunch of other things.
6) See #4
7) See #4, ipv6 version
8) See #4, ipv6 version
```
