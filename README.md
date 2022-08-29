# nft-monitor-service

- [nft-monitor-service](#nft-monitor-service)
    - [Description](#description)
    - [Requirements](#requirements)
    - [Instructions](#instructions)
    - [todos](#todos)
- [Environment variables](#environment-variables)
- [Resources](#resources)

## Description

The service monitors blockchain-Solana and searches for NFT. Next step, the service parses information about nft and
displays it in the console.

## Requirements

- Go 1.17
- Docker
- Gomock

## TODOS

- increase test coverage
- for free api, multi-threaded reading of blocks is [not necessary](!https://docs.alchemy.com/reference/compute-units),
  but if we want to use the business api we need to do this.

## Instructions

### Build

To run build make sure you have docker installed.

```
make build
```

### ReGen Mock

```
mockgen -source=./internal/usecase/interfaces.go -destination=./internal/usecase/mock/interfaces_mock.go -package=mock

mockgen -source=./pkg/solana/solana.go -destination=./pkg/solana/mock/solana_mock.go -package=mock
```

### Running Tests

After that use the command:

```
make test 
```

### Start execution

need set environment variables (see below)

```shell
./monitor
```

### Start with docker

build image

```shell
docker-compose build
```

start container

```shell
docker-compose up -d
```

# Environment variables

| env variable name | description                                                            | default |
|-------------------|------------------------------------------------------------------------|---------|
| LOGGER_LEVEL      | Logging level for logrus framework.                                    | info    |
| START_BLOCK       | block to start parsing from.                                           | None    |
| WAIT_TIME_BLOCK   | how long to wait for the next finalized block(in seconds).             | None    |
| TIMEOUT           | Request timeout.                                                       | None    |
| URL               | node provider url.                                                     | None    |
| API_KEY           | Key for authorization in the node provider.                            | None    |
| API_KEY           | which Solana program should be parsed (Metaplex NFT Candy Machine v2). | None    |

# Resources

| URL                                             | description               |
|-------------------------------------------------|---------------------------|
| https://docs.alchemy.com/reference/api-overview | node provider             |
| https://github.com/evrone/go-clean-template/    | template                  |
| https://solscan.io/                             | block explorer for solana |
| https://docs.solana.com/                        | solana docs               |


