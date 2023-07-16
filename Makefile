.PHONY: run

run:
	go run cmd/main.go & python3 worker.py
