.PHONY: default
default: test


.PHONY: check-working-tree
check-working-tree:
	@git diff-index --quiet HEAD -- \
	|| (echo "Working tree is dirty. Commit all changes."; false)


.PHONY: coverage
coverage: test
	go tool cover -html=dist/coverage.txt


.PHONY: dynamodb
dynamodb:
	@scripts/docker-up-localdev


.PHONY: production
production: check-working-tree
	aws ecr get-login-password \
	| docker login \
	    --username AWS \
	    --password-stdin \
	    `scripts/get-docker-registry production`
	scripts/build-and-deploy production


.PHONY: staging
staging: check-working-tree
	aws ecr get-login-password \
	| docker login \
	    --username AWS \
	    --password-stdin \
	    `scripts/get-docker-registry staging`
	scripts/build-and-deploy staging


.PHONY: run
run:
	docker-compose -f docker-compose.localdev.yml pull --include-deps
	foreman start -e /dev/null


.PHONY: test
test:
	@scripts/docker-up-test
	@echo ' ____'
	@echo '|  _ \ __ _ ___ ___ '
	@echo '| |_) / _` / __/ __|'
	@echo '|  __/ (_| \__ \__ \'
	@echo '|_|   \__,_|___/___/'
