# antha
[![GoDoc](http://godoc.org/github.com/antha-lang/antha?status.svg)](http://godoc.org/github.com/antha-lang/antha)
[![Build Status](https://travis-ci.org/antha-lang/antha.svg?branch=master)](https://travis-ci.org/antha-lang/antha)

Antha v0.0.2

## Installation Instructions

### Docker

Docker is a virtualization technology that allows you to easily download and run
pre-build operating system images on your own machine. 

  1. Install [docker](https://www.docker.com) by following the instructions for
     your operating system.
  2. Follow the docker instructions to start the docker server on your machine
  3. Run the antha docker image
```bash
docker run -it antha/antha
```
  4. Inside the antha docker image, you can follow the instructions below
     for making and running antha elements

By default, when you run the antha image again, you will get a new machine
instance and any changes you made in previously will not be available. If you want
to persist changes you make to antha elements, you can mount directories on your
host machine inside docker instances.

For example,
  1. Download the antha github repo on your (host) machine.
```bash
git clone https://github.com/antha-lang/antha
```
  2. Run the antha docker image and mount your host directory to a directory
     inside the docker instance.
```bash
docker run -it -v `pwd`/antha:/go/src/github.com/antha-lang/antha antha/antha
```
Now, any changes you make to antha on your host machine will be available
within the docker instance.

### OSX (Native)

First step is to install or upgrade to go 1.6. Follow the instructions at the
[Golang](http://golang.org/doc/install) site. 

After you install go, if you don't have [Homebrew](http://brew.sh/), please
install it. Then, follow these steps to setup a working antha development
environment:
```bash
# Setup environment variables
cat<<EOF>>$HOME/.bash_profile
export GOPATH=$HOME/go
export PATH=\$PATH:$HOME/go/bin
EOF

# Reload your profile
. $HOME/.bash_profile

# Install the xcode developer tools
xcode-select --install

# Install some external dependencies
brew update
brew install pkg-config homebrew/science/glpk sqlite3 opencv

# Install antha
go get github.com/antha-lang/antha/cmd/...
```

### Linux (Native)

Depending on your Linux distribution, you may not have the most recent version
of go available from your distribution's package repository. We recommend you
[download](https://golang.org/) go directly. 

For Debian-based distributions like Ubuntu on x86_64 machines, the installation
instructions follow.  If you do not use a Debian based system or if you are not
using an x86_64 machine, you will have to modify these instructions by
replacing the go binary with one that corresponds to your platform and
replacing ``apt-get`` with your package manager.
```bash
# Install go
curl -O https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.6.2.linux-amd64.tar.gz

# Setup environment variables
cat<<EOF>>$HOME/.bash_profile
export GOPATH=$HOME/go
export PATH=\$PATH:/usr/local/go/bin:$HOME/go/bin
EOF

# Reload your profile
. $HOME/.bash_profile

# Install antha external dependencies
sudo apt-get install -y libglpk-dev libopencv-dev libsqlite3-dev git

# Now, we are ready to get antha
go get github.com/antha-lang/antha/cmd/...
```

### Windows (Native)

Installing antha on Windows is significantly more involved than for OSX or
Linux. The basic steps are:

  - Setup a go development environment:
    - Install the source code manager [git](https://git-scm.com/download/win)
    - Install [go](https://golang.org/dl/)
    - Install the compiler [mingw](http://sourceforge.net/projects/mingw/files/Installer/mingw-get-setup.exe/download).
      Depending on whether you installed the 386 (32-bit) or amd64 (64-bit) version
      of go, you need to install the corresponding version of mingw.
  - Download antha external dependencies
    - Install [glpk](http://sourceforge.net/projects/winglpk/) development library and make sure that
      mingw can find it.

If this procedure sounds daunting, you can try using some scripts we developed
to automate the installation procedure on Windows.
[Download](scripts/windows/windows-install.zip), unzip them and run
``install.bat``. This will try to automatically apply the Windows installation
procedure with the default options. Caveat emptor.

## Checking Your Installation

After following the installation instructions for your machine. You can check
if Antha is working properly by running a test protocol
```bash
cd $GOPATH/src/github.com/antha-lang/antha/antha/examples/workflows/constructassembly
antharun --workflow workflow.json --parameters parameters.yml
```

## Making and Running Antha Components

The easiest way to start developing your own antha components is to place them
in the ``antha/component/an`` directory and follow the structure of the
existing components there. Afterwards, you can compile and use your components
with the following commands:
```bash
cd $GOPATH/src/github.com/antha-lang/antha
make
go get github.com/antha-lang/antha/cmd/...
antharun --workflow myworkflowdefinition.json --parameters myparameters.yml
```

## Demo 

[![asciicast](https://asciinema.org/a/12zsgt153sffmfnu2ym7vq9d2.png)](https://asciinema.org/a/12zsgt153sffmfnu2ym7vq9d2)

