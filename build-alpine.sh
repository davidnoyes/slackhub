#!/bin/bash
set -eu
CGO_ENABLED=0 go build -a -installsuffix cgo
