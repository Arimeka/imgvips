.PHONY: lint
lint:
	@scripts/linter.sh

.PHONY: test
test: lint
	@scripts/test.sh
