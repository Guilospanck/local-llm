#!/bin/bash

echo 
echo "Ollama started in the background..."
/bin/ollama serve &

echo 
echo "⏳ Waiting for Ollama to be ready..."
until /bin/ollama list >/dev/null 2>&1; do
	sleep 2
done

echo 
echo "📥 Pulling models inside Ollama container..."
/bin/ollama pull deepseek-r1:1.5b 
/bin/ollama pull gemma:2b
/bin/ollama pull llama3.2

echo 
echo "✅ Model downloaded. Exiting..."
