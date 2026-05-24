BUF_BIN=buf

.PHONY: gen lint clean

gen:
	@echo "Generating proto code..."
	$(BUF_BIN) generate

lint:
	$(BUF_BIN) lint

clean:
	rm -rf gen/go/* gen/ts/*
