#!/bin/sh

curl -s wttr.in/$1 2>/dev/null | sed -r "s/\[([0-9]*;?)+[mGK]//g" | head -n 7
