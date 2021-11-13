#!/bin/bash

. utility.sh

cur_dir=$(pwd)
user=$(whoami)
check_go=$(go version 2> /dev/null)
if [ -z "$check_go" ]; then
    Alert "Go is not installed"
    if [ "$cur_dir" != "/home/${user}" ]; then
        Alert "You are not at directory /home/${user}"
        Alert "Redirecting back to /home/${user}"
        cd /home/$user
    fi
    Alert "Installing Go Now"
    wget https://golang.org/dl/go1.16.6.linux-amd64.tar.gz && tar -C "/home/${user}" -xzf go1.16.6.linux-amd64.tar.gz
    echo "export PATH=${PATH}:/home/${user}/go/bin" >> .profile
    source .profile
    Info "Go is installed. Version : $(go version)"
    Info "Cleaning up..."
    rm go1.16.6.linux-amd64.tar.gz
else
    Info "Go is installed. Version : ${check_go}"
fi

check_git=$(git --version 2> /dev/null)
if [ -z "$check_git" ]; then
    Alert "Git is not installed"
    Alert "Installing Git Now"
    sudo apt update -y 
    sudo apt install git -y pass
    Info "Git is installed. Version : $(git --version)"
else
    Info "Git is installed. Version : ${check_git}"
fi

check_curl=$(curl --version 2> /dev/null)
if [ -z "$check_curl" ]; then
    Alert "cURL is not installed"
    Alert "Installing cURL Now"
    sudo apt update -y 
    sudo apt install curl -y pass
    Info "cURL is installed. Version : $(cURL --version)"
else
    Info "cURL is installed. Version : ${check_curl}"
fi

check_docker=$(docker --version 2> /dev/null)
if [ -z "$check_docker" ]; then
    Alert "Docker is not installed"
    if [ "$cur_dir" != "/home/${user}" ]; then
        Alert "You are not at directory /home/${user}"
        Alert "Redirecting back to /home/${user}"
        cd /home/$user
    fi
    Alert "Installing Docker Now"
    wget https://download.docker.com/linux/static/stable/x86_64/docker-20.10.7.tgz 
    tar xzvf docker-20.10.7.tgz
    sudo cp docker/* /usr/bin/
    sudo dockerd &
    sleep 15
    sudo chmod 777 /var/run/docker.sock
    sudo docker run hello-world
    Info "Successfully Installed Docker. Version : $(sudo docker --version)"
    Info "Adding access right to use docker as non-root user..."
    sudo groupadd docker
    sudo usermod -aG docker $USER
    Info "Completed. Please log out and log back to re-evaluate your group membership"
    Info "Cleaning up..."
    rm docker-20.10.7.tgz
    rm -rf docker
    Info "Completed"
else
    Info "Docker is installed. Version : ${check_docker}"
fi

check_docker_compose=$(docker-compose --version 2> /dev/null)
if [ -z "$check_docker_compose" ]; then
    Alert "Docker-Compose is not installed"
    if [ "$cur_dir" != "/home/${user}" ]; then
        Alert "You are not at directory /home/${user}"
        Alert "Redirecting back to /home/${user}"
        cd /home/$user
    fi
    Alert "Installing Docker-Compose Now"
    sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    Info "Docker-Compose is installed. Version : $(docker-compose --version)"
    Info "Cleaning up..."
    rm "docker-compose-$(uname -s)-$(uname -m)"
else
    Info "Docker-Compose is installed. Version : ${check_docker_compose}"
fi

Info "Pulling latest Hyperledger Fabric Images..."
curl -sSL https://bit.ly/2ysbOFE | bash -s
echo "export PATH=${PATH}:$PWD/fabric-samples/bin" >> $HOME/.profile
source $HOME/.profile

Info "Completed All the dependency Installation. We are good to go!"
