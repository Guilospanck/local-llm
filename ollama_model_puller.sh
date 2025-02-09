#!/bin/bash

# Start Ollama in the background
/bin/ollama serve &

# Wait for Ollama to be ready
sleep 5

# Pull the model
/bin/ollama pull deepseek-r1

