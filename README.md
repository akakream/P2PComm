# ipfsPubSub vs libp2pPubSub Example Implementations

The Kubo's PubSub RPC API will be gradually deprecated. See [ipfs/kubo/issues/9717](https://github.com/ipfs/kubo/issues/9717). This removes the `ipfs pubsub` commands from IPFS. However, there is still a big community that uses `ipfs pubsub`. This repository contains the same pubsub logic implemented using the `ipfs pubsub` API and the `go-libp2p-pubsub` API to help developers who want to see a comparison between the two APIs.

## Build

```
make build
```

## Run

### libp2p

Start with

```
./bin/app -uselibp2p
```

### ipfs

Start ipfs daemon with pubsub enabled in another terminal to run the ipfs pubsub.

```
ipfs daemon --enable-pubsub-experiment
```

Then start with

```
./bin/app
```
