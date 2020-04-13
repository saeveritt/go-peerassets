#!/bin/bash
set -e

if [ "$(echo "$1" | cut -c1)" = "-" ]; then
  echo "$0: assuming arguments for peercoind"
  set -- peercoind "$@"
fi

if [ "$(echo "$1" | cut -c1)" = "-" ] || [ "$1" = "peercoind" ]; then

  mkdir -p "$PPC_DATA"
  chmod 700 "$PPC_DATA"
  chown -R peercoin "$PPC_DATA"

	if [[ ! -s "$PPC_DATA/peercoin.conf" ]]; then
    cat <<-EOF > "$PPC_DATA/peercoin.conf"
    rpcallowip=::/0
    rpcpassword=${RPC_PASSWORD}
    rpcuser=${RPC_USER}
		EOF
    cat "$PPC_DATA/peercoin.conf"
	fi

  set -- "$@" -datadir="$PPC_DATA"
fi

if [ "$1" = "peercoind" ] || [ "$1" = "peercoin-cli" ] || [ "$1" = "peercoin-tx" ]; then
  echo
  exec su-exec peercoin "$@"
fi

exec "$@"