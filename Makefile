obu:
	@go build -o ./bin/obu ./obu
	@./bin/obu

receiver:
	@go build -o ./bin/receiver ./receiver
	@./bin/receiver

calculator:
	@go build -o ./bin/calculator ./distance_calculator
	@./bin/calculator

aggregator:
	@go build -o ./bin/aggregator ./aggregator
	@./bin/aggregator

.PHONY: obu receiver calculator aggregator