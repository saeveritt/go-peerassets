#!/bin/sh
set -e
curl -X POST -d "blockhash=$1" "http://go-peerassets:8089/v1/alert"