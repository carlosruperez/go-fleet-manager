BINARIES_DIRECTORY = bin
BASIC_TAG = basic
MANAGER_TAG = manager
BASIC_DIRECTORY = $(BINARIES_DIRECTORY)/basic
MANAGER_DIRECTORY = $(BINARIES_DIRECTORY)/manager

BASE_BUILD_BASIC_COMMAND = go build --tags $(BASIC_TAG),main -ldflags "-w -s" -o $(BASIC_DIRECTORY)/go-fleet-$(BASIC_TAG)-$@
BASE_BUILD_MANAGER_COMMAND = go build --tags $(MANAGER_TAG),main -ldflags "-w -s" -o $(MANAGER_DIRECTORY)/go-fleet-$(MANAGER_TAG)-$@

build: windows_amd64 linux_amd64 macos_arm64

windows_amd64:
	GOOS=windows GOARCH=amd64 $(BASE_BUILD_BASIC_COMMAND).exe
	GOOS=windows GOARCH=amd64 $(BASE_BUILD_MANAGER_COMMAND).exe

linux_amd64:
	GOOS=linux GOARCH=amd64 $(BASE_BUILD_BASIC_COMMAND)
	GOOS=linux GOARCH=amd64 $(BASE_BUILD_MANAGER_COMMAND)

macos_arm64:
	GOOS=darwin GOARCH=arm64 $(BASE_BUILD_BASIC_COMMAND)
	GOOS=darwin GOARCH=arm64 $(BASE_BUILD_MANAGER_COMMAND)

run-basic:
	go run --tags $(BASIC_TAG) . $(ARGS)

run-manager:
	go run --tags $(MANAGER_TAG) . $(ARGS)
