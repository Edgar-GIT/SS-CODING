#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
LIB="$(pwd)/discord/music_bot/libdave"
export CGO_ENABLED=1
export CGO_CFLAGS="-I${LIB}/include"
export CGO_LDFLAGS="-L${LIB}/lib -Wl,-rpath,${LIB}/lib -ldave -lstdc++"
if [ ! -f "${LIB}/include/dave/dave.h" ]; then
  go run ./tools/bootstrap
fi
exec go run ./manager "$@"
