#!/bin/bash

cmd=$1
if [[ $cmd == 'download' ]]
then
    go run cmd/bdpan/bdpan.go -d ~/Downloads $@
else
    go run cmd/bdpan/bdpan.go $@
fi
