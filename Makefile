.PHONY: test

.DEFAULT: test

test:
	docker build -t test .
