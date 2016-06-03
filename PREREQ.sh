#!/bin/sh

# install_protoc installs protobuf by getting and compiling its source code.
install_protoc() {
    echo "Installing the latest version of protocol buffer..."
    sudo apt-get install build-essential
    sudo apt-get install autoconf
    sudo apt-get install libtool
    git clone git@github.com:google/protobuf.git $GOPATH/src/github.com/google/protobuf
    cd $GOPATH/src/github.com/google/protobuf
    ./autogen.sh
    ./configure
    make
    make check
    sudo make install
    sudo ldconfig
}

# Installing protobuf if not install, or it is not version 3.0+.
if ! hash protoc 2>/dev/null; then
    install_protoc
else
    version=`protoc --version | awk '{printf("%d.%d", $2, $3);};'`
    if [ $(echo $version' < 3.0' | bc) -eq "1" ]; then
	install_protoc
    else
	echo "Protocol buffer is already the latest version."
    fi
fi

# Installing lit package dependencies
echo "Installing lit package dependencies..."
go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u golang.org/x/net/context
echo "Done!"
