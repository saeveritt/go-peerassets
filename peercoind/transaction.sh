#!/bin/sh
set -e
curl -X POST -d "txid=$1" "http://localhost:8089/alert"