# Run Node Devnet infrastrucuture

## Generate Devnet state

First you need to generate Anvil state executing the following commands in a shell 
```
cd anvil
docker compose up
```

Wait until you see a log like the one bellow and stop execution 
```
anvil-hardhat-1 exited with code 0
```

## Run infrastructure with sample echo DApp

After the base Anvil state files are generated (see section above), at `devnet` folder, run :
```
docker compose up
```

It will run a very simple Echo Dapp that generates an notice with the same payload sent as an input

## Send input

Using `cast` , to send an input run :

```
INPUT=0x68656C6C6F206E6F6465 ;\
INPUT_BOX_ADDRESS=0x59b22D57D4f067708AB0c00552767405926dc768 ;\
DAPP_ADDRESS=0x70ac08179605AF2D9e75782b8DEcDD3c22aA4D0C ;\
cast send $INPUT_BOX_ADDRESS "addInput(address,bytes)(bytes32)" $DAPP_ADDRESS $INPUT --mnemonic "test test test test test test test test test test test junk" --rpc-url "http://localhost:8545"
```

## Query input

Execute the following query to get all the produced notices and the associated input
```
query{
  notices {edges{node{index,payload,input{payload}}}}
}
```
Executing this query using `curl`
```
curl 'http://localhost:4000/graphql' -H 'Accept-Encoding: gzip, deflate, br' -H 'Content-Type: application/json' -H 'Accept: application/json' -H 'Connection: keep-alive' -H 'DNT: 1' -H 'Origin: http://localhost:4000' --data-binary '{"query":"query{  notices {edges{node{index,payload,input{payload}}}}}"}' --compressed
```