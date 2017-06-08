VER=3

docker: fe
	docker build -t jeffallen/emojch-fe:$(VER) .
	docker push jeffallen/emojch-fe:$(VER)

fe: *.go
	go-bindata-assetfs -ignore '~$$' static/...
	CGO_ENABLED=0 go1.8 build -a -installsuffix cgo
