# A Backend For NFT Marketing With Privacy Security And Supervision

This project is based on HyperLedger Fabric.It implementes a backend, but its not a completely availabel system.

## Table of Contents

- [Security](#security)
- [Install](#install)
- [Usage](#usage)
- [API](#api)
- [License](#license)

## Security

* Supervision Model
* NFT Lifecycle Encryption Based On SM2
* User Operations Encryption Based On SM2
* User Register Encryption (Not Implemented)
* Cert Verify Method

## Install

Make sure you have installed go >= 1.18.3 and docker >= 20.10.17.You should use  `git clone` to add this project  locally and cd into the folder.Run ./bootstrap.sh to install fabric.

```bash
git clone https://github.com/xiahezzz/NFTChain.git && cd /NFTChain/
./bootstrap.sh
```

## Usage

Now,you can start the blockchain,and deploy this chaincode.

```bash
cd /blockchain/network
make networkup
make createChannel
make deployCC
```

Before you start the service,you shou download the go module ,which used in the service.

```bash
cd ../../service
go mod download && go mod verify
```

Congratulations！You can build the app and run it.

```
go build -o app main.go
./app
```

At last,stop your service,stop and clean the network.

```
cd ../blockchain/network
make networkdown
```

## API

APIs in the service are shown in SetupRouter() function.

## License

[MIT © xiahezzz.](../LICENSE)
