#!/bin/bash

# from:
# https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html#docker

# Setting up NVIDIA Container Toolkit
# Setup the package repository and the GPG key
distribution=$(. /etc/os-release;echo $ID$VERSION_ID) \
      && curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg \
      && curl -s -L https://nvidia.github.io/libnvidia-container/$distribution/libnvidia-container.list | \
            sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' | \
            sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list

# install nvidia-docker2
sudo apt-get update
sudo apt-get install -y nvidia-docker2

# restart docker
sudo systemctl restart docker

# test
# sudo docker run --rm --gpus all nvidia/cuda:11.0.3-base-ubuntu20.04 nvidia-smi