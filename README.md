# Shurli

 Shurli App - alpha version

## :warning: IMPORTANT NOTE: Please use only small amounts of funds in your wallet while testing subatomic swaps

## :warning: This is in alpha phase right now so please expect breaking changes.


## Objective

 To build an independent GUI application that works with existing supported wallets to make sub-atomic trades.

### Currently supported

* Komodod
* Komodo-ocean-qt
* Verus-desktop

### Code of conduct

* While reporting issues, please report all the debug data at [Shurli Issues](https://github.com/Meshbits/shurli/issues).

### Requirements

    - Go v1.14+
    - Git
    - Komodo Daemon (dev branch compiled from http://github.com/jl777/komodo.git)

### Install Git, Komodo and Subatomic binaries

Follow this instruction guide and install `Komodo` and `subatomic` binaries on your system:
https://gist.github.com/himu007/add3181427bb53ab5dc5160537f0c238

#### Windows Komodo and Subatomic binaries
At this stage, the required subatomic related API is only available via jl777's copy of komodo source code. So, you can either cross-compile the windows binaries and use those on your windows machine or you can try the debug komodo binaries we used in our development environment, which are available here:
https://github.com/Meshbits/komodo/releases/tag/debug_release

**It's best to not use these binaries for big amount of funds. Just use the linked `komodod` to start `DEX` blockchain which provides the required DEXP2P API. You can use other GUI wallets as usual for other subatomic supported cryptocurrencies.**

### Install Go on your system

On Ubuntu can follow this guide: https://github.com/golang/go/wiki/Ubuntu

```shell
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go
```

On Mac can install Go using `brew`.

```shell
brew update
brew install go
```

### Installing Shurli App

In Linux/Mac you must have a `go` directory under your `$HOME` directory in your OS.
If you don't find this directory then create the following:

```shell
mkdir -p $HOME/go/{bin,src,pkg}
```

```
go get -u github.com/Meshbits/shurli
```

#### Copy or symlink Komodo and Subatomic binary to Shurli

Once you have the `komodod`, `komodo-cli` and `subatomic` executables available (either by compilation or downloaded), you need to have these binaries accessible via Shurli's `assets` directory.
You'll find `assets` directory within the Shurli files.

The directory structure of the Shurli looks like this:

```shell
.
├── assets
│   └── subatomic.json
├── chains.json
├── config.json.sample
├── favicon.png
├── main.go
├── public
│   ├── coins
│   ├── css
│   ├── fapro
│   ├── gfonts
│   ├── images
│   └── js
├── sagoutil
├── templates
```

Just copy (or symlink on Linux/Windows) the `komodod`, `komodo-cli` and `subatomic` binaries to the `assets/` directory where you see `subatomic.json` file located.

#### Configure config.json

Make a copy of `config.json` file from `config.json.sample` file:

```shell
cd $HOME/go/src/github.com/Meshbits/shurli
cp config.json.sample config.json
```

<!--
Open `config.json` in text editor and edit value of only `subatomic_dir` key.

To get the full path of your `komodo` directory `cd` to `komodo` directory and issue `pwd` command to get full path. Example:

```shell
cd $HOME/komodo/src
pwd
```

It will give full path for example like this:

```shell
/home/satinder/komodo/src
```

Replace the path you get from this output in `config.json` file.

Before:

```json
"subatomic_dir": "/Users/satinder/repositories/jl777/komodo/src"
```

After:

```json
"subatomic_dir": "/home/satinder/komodo/src"
```

For Windows, make sure to use the following format for setting up `subatomic_dir` path:
```json
"subatomic_dir": "C:/Users/satinder/kmdsub"
```

Note that it's not backslash `\` but forward slash `/` for the path.
If the format of this path would be incorrect, Shurli will have issue locating the `subatomic` binary on your machine.
-->

#### Configure DEX blockchain's parameters in `config.json` file

You MUST update value of `dex_pubkey`, `dex_handle`, `dex_recvzaddr`, `dex_recvtaddr` in your **config.json** file before starting Shurli application.

- `dex_recvtaddr`: is your KMD's public address. The address which starts with letter `R`.
- `dex_pubkey`: is pubkey of your KMD's public key. You can it from the your KMD wallet application.
- `dex_recvzaddr`: is the PIRATE's private address. The address which starts with letter `z`.
- `dex_handle`: it is very much like a unique username you want to use on Subatomic swaps. It will show to other traders when they will see your orders in orderbook. Your handle must with without space between letters.

```json
    "dex_nSPV": "1",
    "dex_addnode": "136.243.58.134",
    "dex_pubkey": "03_YOUR_PUBKEY_FROM_DEX",
    "dex_handle": "SET_YOUR_HANDLE_FOR_SUBATOMIC",
    "dex_recvzaddr": "YOUR_PIRATE_PRIVATE_ADDRRESS",
    "dex_recvtaddr": "YOUR_KMD_PUBLIC_ADDRRESS"
```

#### Build Shurli application

```shell
cd $HOME/go/src/github.com/Meshbits/shurli
make
```

#### Start Shurli App

The build command will make a system executable binary named "shurli".
To start Shurli execute the following command:

```shell
./shurli start
```

It will start Shurli in daemon mode, leaving Shurli running in background.

Now open http://localhost:8080

#### Stop Shurli App

To stop Shurli you have to execute the stop command.
Otherwise it will keep running in background.
Stop with following command:

```shell
cd $HOME/go/src/github.com/Meshbits/shurli
./shurli stop
```

#### Shurli logs

If starting Shurli as daemon using `./shurli start` command, it will not show any logs on cosole output.
To view the logs you can check `shurli.log` file in Shurli directory.
Following example command on Linux/OSX will show updated prints being pushed to `shurli.log` file:

```shell
cd $HOME/go/src/github.com/Meshbits/shurli
tail -f shurli.log
```

And Windows users can use the following command in PowerShell to check live shurli logs:
```shell
Get-Content .\shurli.log -Wait
```

you can press CTRL+C to cancel `tail` or `Get-Content` command's output.

#### Making a release build

Release builds can be made cross platform.
Means you can build Mac OS build on Linux, and Linux builds on Mac OS,
thanks to Go's cross-compilation capabilities.

##### Linux build

To make Linux distributable build execute the following command:

```shell
cd $HOME/go/src/github.com/Meshbits/shurli
make build-linux
```

After this command you'll find a directory `dist/dist_unix` in `$HOME/go/src/github.com/Meshbits/shurli/`.
All the required files for Shurli would be in `dist/dist_unix`. You can renamed or moved `dist_unix` to anywhere on the machine.
Or make a zip or tar.gz archive of it to distribute.

##### Mac OS build

To make Mac OS distributable build execute the following command:

```shell
cd $HOME/go/src/github.com/Meshbits/shurli
make build-osx
```

Smiliar to Linux build, for Mac OS you'll find `dist/dist_osx` in `$HOME/go/src/github.com/Meshbits/shurli/`.
You can rename, move or archive the `dist_osx` and distribute.

##### Windows build

To make Windows distributable build execute the following command:

```shell
cd %USERPROFILE%\go\src\github.com\Meshbits\shurli
make build-win
```

You'll find the windows build files in `dist/dist_win`.


##### Clean build

To clean all compiled files execute the following command:

```shell
cd $HOME/go/src/github.com/Meshbits/shurli
make clean
```

It will delete all dist and binary files the build commands created.
