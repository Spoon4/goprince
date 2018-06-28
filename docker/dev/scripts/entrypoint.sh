#!/bin/ash

make --no-print-directory run
dlv debug --headless --listen=:2345 --log=true --api-version=2 goprince

exec "$@"