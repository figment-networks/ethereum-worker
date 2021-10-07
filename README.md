# ethereum-worker

## Worker

Stateless worker is responsible for connecting with the chain and retrieving account balance and total token balance for erc20 chains by accountAddress, contractAddress or Network name

### Compile

To compile sources you need to have go 1.14.1+ installed.

```bash
    make build-live
```

### Running

```
docker run -p 8097:8097 -it eth-worker
```

In order to run, it requires a tunnel to the ethereum--archive-1 on port 8545. Once this is done you can access data by endpoints such as

```
http://localhost:8097/getBalance?accountAddress=0x9320e85de19928f60387be5ac553791bebcdf2d3&contractAddress=0x00c83aeCC790e8a4453e5dD3B0B4b3680501a7A7

http://localhost:8097/getTotalSupply?network=skale

```
