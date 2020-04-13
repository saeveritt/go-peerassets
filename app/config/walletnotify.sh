#!/bin/sh
set -e
curl -X POST -d "txid=$1" "http://go-peerassets:8089/v1/alert"