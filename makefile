.PHONY: default
default: test

.PHONY: coverage
coverage: test
	go tool cover -html=dist/coverage.txt

.PHONY: test
test:
	@scripts/docker-up-test
	@echo ' ____'
	@echo '|  _ \ __ _ ___ ___ '
	@echo '| |_) / _` / __/ __|'
	@echo '|  __/ (_| \__ \__ \'
	@echo '|_|   \__,_|___/___/'
