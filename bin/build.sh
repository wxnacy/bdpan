#!/bin/bash

# cd cmd/main && go build -o bdpan && mv bdpan $(go env GOPATH)/bin && cd --
cd cmd/bdpan && go build && mv bdpan $(go env GOPATH)/bin && cd --
