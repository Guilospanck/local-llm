#!/bin/bash

back_watch() {
	cd back/ && go mod tidy && air
}

back() {
	cd back/ && go mod tidy && go run .
}

front() {
	cd front/ && pnpm i && pnpm dev
}

if [ "$1" == "--watch" ]; then
  echo "Watch mode enabled"
  front & back_watch
else
  echo "Watch mode disabled"
  front & back
fi
