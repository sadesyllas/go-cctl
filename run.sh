#!/bin/bash

function kill_cctl {
  killall go-cctl
}

trap kill_cctl INT

trap kill_cctl TERM

trap kill_cctl HUP

cd "$(dirname "$0")"

kill_cctl &> /dev/null

go run . -p 3003 &

cd web

pnpm preview -- --host 0.0.0.0 --port 3005
