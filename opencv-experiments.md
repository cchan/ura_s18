# RaspiVideo / OpenCV Experiments

## Basic raspivid display

To test out video streaming...

- [if you're on Windows Bash, make sure you have XMing set up]
- `raspberry pi $ raspivid -hf -w 400 -h 400 -t 999999999 -fps 30 -b 5000000 -o - | nc -l -p 8160`
- `computer $ sudo apt install mplayer`
- `computer $ nc [RPI IP ADDRESS] 8160 | mplayer -fps 200 -demuxer h264es -` (Note that the 200fps is so that mplayer consumes the stream as fast as it's created, so that it's never delayed. If there's network buffering due to latency, the 200fps quickly consumes the buildup of frames.)

You should get a nice window displaying your live video. :)

Notes:

- The IP address can be found by looking at `ifconfig` and checking the relevant interface - in my case, I'm connecting to RPi over USB, so interface `usb0`.
- You can also just ssh tunnel it since you already have ssh access: `computer $ ssh -R 8160:localhost:8160 pi@raspberrypi.local` (or just use PuTTY)
- To be clear, the port `8160` is entirely arbitrary.

## Python/OpenCV video capture to output

- `sudo modprobe bcm2835-v4l2` to make it visible to opencv (assuming you've enabled it in raspi-config)
- `python3 vidcap.py | nc -l -p 8160`
