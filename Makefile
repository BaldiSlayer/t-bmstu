.PHONY: run

run:
	go run cmd/main.go & python3 workers/timus/timus.py

