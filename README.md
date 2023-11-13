# core
Official Go implementation of the JMESWorld protocol


## Build and run locally

-  `cd cmd/jmesd && go build` - Will build `jmesd` daemon.

## Build dockerfile

` docker buildx build  --build-arg TOOL_NODE_FLAGS="--max-old-space-size=8192 --max_semi_space_size=256" --platform linux/arm --push -t jmesworld-core .`
## Run

Run `/cmd/jmesd/jmesd` with the following options:

- `--chain-id` - Chain ID of the network to connect to