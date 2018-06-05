clean:
	rm -f dbachaind dbachaincli

dbaBin:clean
	cd ../../cmd/dbachaind && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build && cp -f dbachaind ../../DOCKER/dba/
	cd ../../cmd/dbachaincli && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build && cp -f dbachaincli ../../DOCKER/dba/

build:dbaBin
	docker build -t $(dockerImage) .

.PHONY:clean prepare build build_develop
