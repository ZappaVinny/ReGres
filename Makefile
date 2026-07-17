BINARY_NAME=rgs

install:
	@cd rgs-cli && go install ./cmd/rgs
	@echo "rgs cli has been installed use 'rgs' command to run it"
build:
	@cd rgs-cli && go build -o ../$(BINARY_NAME) ./cmd/rgs
	@echo "rgs cli binary has been built use './$(BINARY_NAME)' command to run it"