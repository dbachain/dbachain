clean:
	@rm -f gaiad

prepare:clean
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ../../../github.com/cosmos/cosmos-sdk/cmd/gaia/cmd/gaiad && \
		cp -f gaiad ../../DOCKER/dba/

build:prepare
	docker build -t $(dockerImage) .

.PHONY:clean prepare build_image
