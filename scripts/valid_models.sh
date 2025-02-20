#!/bin/bash

valid_models="deepseek-r1:1.5b gemma:2b llama3.2"

validate_model() {
	if ! echo "$valid_models" | grep -wq "$1"; then
		echo "Error: Invalid model '$1'. Choose from: $valid_models" >&2
		exit 1
	fi
}
