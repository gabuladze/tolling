obu:
	@go build -o ./bin/obu ./obu
	@./bin/obu

receiver:
	@go build -o ./bin/receiver ./receiver
	@./bin/receiver

.PHONY: obu receiver