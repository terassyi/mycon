#! /bin/sh

sudo apt update

# install golang
wget https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.14.2.linux-amd64.tar.gz
rm go1.14.2.linux-amd64.tar.gz

echo "export PATH=$PATH:/usr/local/go/bin" >> .bashrc

# install docker
sudo apt -y install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"

sudo apt update
sudo apt -y install docker.io

sudo systemctl start docker

# add docker user group
sudo groupadd docker
sudo gpasswd -a $USER docker
sudo systemctl enable docker