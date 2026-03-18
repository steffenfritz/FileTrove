#!/bin/sh
set -e

# Ensure /usr/local/lib is in the dynamic linker search path so libyara_x_capi.so is found
echo "/usr/local/lib" > /etc/ld.so.conf.d/ftrove.conf
ldconfig
