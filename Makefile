TEST_FILE=go list ./... | grep -v /contracts
test:
	go test $$($(TEST_FILE)) -coverprofile out/cover.out
coverage: clean test
	go tool cover -html=out/cover.out
clean:
	@echo "cleaning out folder..."
	@rm -rf out/*