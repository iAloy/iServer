# **iServer**
Basic fileserver written in go

This is a basic file hosting server for accessing and downloading files from a directory.
It also support ignore file which follows mainly gitignore rules (except sub/multilevel gitignore file).

## Usage

```sh
./iServer -p <port> -d <directory_name>
```
Default `port` is `8080` and default `directory` is the directory from which the programm is called.

All available options are:
- `-i` - ignore file (default is ".ignore")
- `-p` - define port (default is "8080")
- `-d` - directory (default is ".")
- `-h` - show help message
- `-v` - verbose output
- `-V` - show version

## Install Golang in Rpi

https://akashrajpurohit.com/blog/installing-the-latest-version-of-golang-on-your-raspberry-pi/

```sh
#!/bin/bash

VERSION=1.21.1

## Download the latest version of Golang
echo "Downloading Go $VERSION"
wget https://dl.google.com/go/go$VERSION.linux-armv6l.tar.gz
echo "Downloading Go $VERSION completed"

## Extract the archive
echo "Extracting..."
tar -C ~/.local/share -xzf go$VERSION.linux-armv6l.tar.gz
echo "Extraction complete"

## Detect the user's shell and add the appropriate path variables
SHELL_TYPE=$(basename "$SHELL")

if [ "$SHELL_TYPE" = "zsh" ]; then
    echo "Found ZSH shell"
    SHELL_RC="$HOME/.zshrc"
elif [ "$SHELL_TYPE" = "bash" ]; then
    echo "Found Bash shell"
    SHELL_RC="$HOME/.bashrc"
elif [ "$SHELL_TYPE" = "fish" ]; then
    echo "Found Fish shell"
    SHELL_RC="$HOME/.config/fish/config.fish"
else
    echo "Unsupported shell: $SHELL_TYPE"
    exit 1
fi

echo 'export GOPATH=$HOME/.local/share/go' >> "$SHELL_RC"
echo 'export PATH=$HOME/.local/share/go/bin:$PATH' >> "$SHELL_RC"

## Verify the installation
if [ -x "$(command -v go)" ]; then
    INSTALLED_VERSION=$(go version | awk '{print $3}')
    if [ "$INSTALLED_VERSION" == "go$VERSION" ]; then
        echo "Go $VERSION is installed successfully."
    else
        echo "Installed Go version ($INSTALLED_VERSION) doesn't match the expected version (go$VERSION)."
    fi
else
    echo "Go is not found in the PATH. Make sure to add Go's bin directory to your PATH."
fi

## Clean up
rm go$VERSION.linux-armv6l.tar.gz
```

## Building

To build the executable you need to install go compiler and make first.
Then clone the repository and run in the repository
```sh
make
```
Or to build manually run
```sh
go build
```
However for optimized binary size and cross compiling for amd64 and arm64 architecture use of make to build is recommended.

## Roadmap

- [x] add basic file servering
- [x] add directory listing
- [x] support for custom port and folder
- [x] add support for ignore file
- [x] add column for file modification date
- [ ] add download button for each file
- [ ] add support for sorting files and directories
- [ ] add service file for linux
- [ ] make deb package




