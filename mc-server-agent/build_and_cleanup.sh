#!/bin/bash

set -e 

install_go() {
    echo "Installing Go..."

    wget https://go.dev/dl/go1.21.1.linux-amd64.tar.gz -O /tmp/go.tar.gz

    sudo tar -C /usr/local -xzf /tmp/go.tar.gz

    export PATH=$PATH:/usr/local/go/bin
    export GOROOT=/usr/local/go
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    echo "Go installed successfully!"
}

remove_go() {
     echo "Removing Go..."

    sudo rm -rf /usr/local/go

    rm -f /tmp/go.tar.gz

    echo "Go removed successfully!"
}

build_go_binary() {
    echo "Building Go binary..."

    cd .

    go build -o /usr/local/bin/agent

    chmod +x /usr/local/bin/agent

    echo "Go binary compiled and placed in /usr/local/bin!"
}

install_go
build_go_binary
remove_go

echo "Script finished!"
