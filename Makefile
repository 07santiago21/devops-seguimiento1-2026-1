BINARY   = bootstrap
DIST_DIR = dist
ZIP      = $(DIST_DIR)/function.zip

.PHONY: build clean


build:
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(DIST_DIR)/$(BINARY) .
	cd $(DIST_DIR) && zip function.zip $(BINARY)
	@echo "Artifact ready: $(ZIP)"

## Remove build artifacts.
clean:
	rm -rf $(DIST_DIR)
