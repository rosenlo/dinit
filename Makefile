include $(realpath $(GOPATH)/src/dinit/build/Makefile)

.PHONY:default
default: dinit
	#@cd $(WORKSPACE)/build && bash ./build.sh

.PHONY:test
test:
		$(GOTEST) -v ./...

.PHONY:clean
clean:
		$(GOCLEAN)
		@rm -rf $(BIN_PATH)

dinit:
	cd $(WORKSPACE)/cmd/dinit/ && make VERSION=${VERSION}
	docker build -t dinit:${VERSION}\
		-f build/Dockerfile-dinit ./build

push:
	docker push dinit:latest

latest: default
	docker build -t dinit:latest \
		-f build/Dockerfile ./build

# all: default image push
