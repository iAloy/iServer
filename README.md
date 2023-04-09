# iServer 
Basic fileserver written in go

This is a basic file hosting server for accessing and downloading files from a directory.

## Usage

```sh
./iServer -p <port> -d <directory_name>
```
Default `port` is `8080` and default `directory` is the directory from which the programm is called.

All available options are:
- `-p` - define port (default is "8080")
- `-d` - directory (default is ".")
- `-h` - show help message
- `-v` - verbose output
- `-V` - show version 

## Building

To build the executable you need to install go compiler and make first.
Then clone the repository and run in the repository 
```sh
make
```
Or to build manually run 
```sh
go build main.go
```
However for optimized binary size and cross compiling for amd64 and arm64 architecture use of make to build is recommanded.

## Roadmap

- [x] add basic file servering
- [x] add directory listing
- [x] support for custom port and folder
- [ ] add support for ignore file
- [ ] add column for file modification date
- [ ] add download button for each file
- [ ] add support for sorting files and directories




