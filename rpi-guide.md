# Raspberry Pi Guide

Clive Chan, June 2018

This assumes RPi 3B or RPi Zero W, and Raspbian Stretch. It might work on older RPis but I don't know. It almost definitely doesn't work for previous versions of Raspbian.

## Flashing Raspbian and setting up headless
Flash the SD card with the latest version of Raspbian, following [these instructions](https://www.raspberrypi.org/documentation/installation/installing-images/).

Create two files in the boot partition: `ssh`, which is blank, and `wpa_supplicant.conf`, which contains:

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

It's good practice to do this set up process on secured wifi (I like to do it on my phone hotspot), but if you really need to you can remove the `psk=...` line and use `key_mgmt=NONE` for WEP wifi. Don't ask about WPA-Enterprise (i.e. eduroam); it's really hard and eduroam also additionally refuses to let you talk to other computers on the same network, for obvious security reasons.

Remove the SD card, put it into the Pi, and plug it in to usb power. Wait a bit for it to do a first-time boot-up. Then to connect to it, make sure you're on the same network as provided in `wpa_supplicant.conf` and run `ssh pi@raspberrypi`. It'll probably tell you that the authenticity of host 'raspberrypi' cannot be established, which is correct and is one of the reasons you should be doing this setup on secured wifi with no other users. The password is `raspberry`.

Upon logging in, set the password with `passwd` immediately. Please.

It's also generally a good idea to run `sudo apt update` then `sudo apt dist-upgrade` on any new system to update all your packages. If you're wondering, [this](https://askubuntu.com/a/226213) is a good answer describing the difference between `upgrade` and `dist-upgrade`. It's probably a bad idea to do this on cell data, so see the Wifi section for how to do wifi.

## Wifi

Add more networks by editing `/etc/wpa_supplicant/wpa_supplicant.conf`, or you can use `sudo raspi-config` for an even easier to use interface.

If you have lots of networks you can set `priority` for each one. By default it's `0`, and higher priority means it gets selected first. Negative is allowed.

## Protips

Use `sudo shutdown now` before unplugging your Pi to avoid bad things that usually don't happen but occasionally do.

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

for security. You can now safely access your raspberry pi on whatever networks you want.

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
