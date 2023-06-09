# ipfsPubSub vs libp2pPubSub Example Implementations

The Kubo's PubSub RPC API will be gradually deprecated. See [ipfs/kubo/issues/9717](https://github.com/ipfs/kubo/issues/9717). This removes the `ipfs pubsub` commands from IPFS. However, there is still a big community that uses `ipfs pubsub`. This repository contains the same pubsub logic implemented using the `ipfs pubsub` API and the `go-libp2p-pubsub` API to help developers who want to see a comparison between the two APIs.

## Build

```
make build
```

## Run

```
./bin/app server --port=3000 --servertype=libp2p --datastore
```

Start ipfs daemon with pubsub enabled in another terminal to run the ipfs pubsub.

```
ipfs daemon --enable-pubsub-experiment
```

## Test

```
make test
```

## P2P client endpoints

| Endpoint  | Description                  |
| --------- | ---------------------------- |
| `/pub`    | Publish a message to a topic |
| `/sub`    | Subscribe to a topic         |
| `/unsub`  | Unsubscribe from a topic     |
| `/topics` | Return all susbcribed topics |

## Datastore endpoints

| Endpoint | Description                           |
| -------- | ------------------------------------- |
| `/list`  | List the keys stored in the datastore |
| `/get`   | Get the key stored in the datastore   |
| `/put`   | Put a key-value pair in the datastore |

## Other endpoints

| Endpoint  | Description                       |
| --------- | --------------------------------- |
| `/health` | Get health information for server |

## TODO

- Create a PutHook for keys: [https://github.com/ipfs/go-ds-crdt/issues/178](https://github.com/ipfs/go-ds-crdt/issues/178) This can be useful.

- Fix the `/get` key calls. Remove the body and add as parameter.
