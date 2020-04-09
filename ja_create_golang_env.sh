#!/bin/bash

export PATH=$GOROOT/bin:$PATH

GOPATH=$(go env GOPATH)
echo "[Jason]: Checking:"
echo "------------"
echo $GOPATH
echo "------------"
# proto-gen- under $GOPATH/bin  ($GOBIN)
export PATH=$PATH:$GOPATH/bin
