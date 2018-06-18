So https://www.pyimagesearch.com/2017/09/04/raspbian-stretch-install-opencv-3-python-on-your-raspberry-pi/ is extremely slow.

use putty to local ssh tunnel from 2222 local to localhost:22 remote
then ssh to localhost -p 2222

https://gist.github.com/hrshovon/70612b719bfda0becde46f0b9d2dfa36/511333c2d31835a7a39ce11d7b263e78c21016cd
(fuse doesn't work on WSL)
(probably will have some detections incorrect :/ )

cmake -D CMAKE_BUILD_TYPE=RELEASE -D OPENCV_EXTRA_MODULES_PATH=../../opencv_contrib-3.4.1/modules -D ENABLE_NEON=ON -D ENABLE_VFPV3=ON -D BUILD_TESTS=OFF -D INSTALL_PYTHON_EXAMPLES=OFF -D BUILD_EXAMPLES=OFF -D CMAKE_TOOLCHAIN_FILE=../platforms/linux/arm-gnueabi.toolchain.cmake -D WITH_TBB=ON -D BUILD_TBB=ON -D BUILD_opencv_python2=ON -D PYTHON2_LIBRARIES=../../libpython2.7.so -D PYTHON2_INCLUDE_PATH=/usr/include/python2.7 -D PYTHON2_NUMPY_INCLUDE_DIRS=/usr/local/lib/python2.7/dist-packages/numpy/core/include -D BUILD_TESTS=OFF -D INSTALL_PYTHON_EXAMPLES=OFF -D BUILD_opencv_python2=ON -D PYTHON2_LIBRARIES=../../libpython2.7.so -D PYTHON2_INCLUDE_PATH=/usr/include/python2.7 -D PYTHON2_NUMPY_INCLUDE_DIRS=/usr/local/lib/python2.7/dist-packages/numpy/core/include -D BUILD_opencv_python3=ON -D PYTHON3_LIBRARIES=../../libpython3.5m.so -D PYTHON3_INCLUDE_PATH=/usr/include/python3.5 -D PYTHON3_NUMPY_INCLUDE_DIRS=/usr/lib/python3/dist-packages/numpy/core/include -D PYTHON2_NUMPY_VERSION=1.13.3 -D PYTHON3_NUMPY_VERSION=1.14.2 ..
(only change is ~ directory paths and adding numpy versions)

make -j8
sudo apt-get install libtbb-dev zlib1g-dev???
pip install protobuf???
maybe just restart make when it fails??? apparently powers past the error on second try

it fails on target all, but builds cv2.so anyway. Huh.

scp the whole build/lib directory to the raspberry pi, and put the directory (e.g. ~/cv2) into the LD_LIBRARY_PATH.


