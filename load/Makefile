VER=1

docker: load
	docker build -t jeffallen/emojch-load:$(VER) .
	docker push jeffallen/emojch-load:$(VER)

load: *.go
	CGO_ENABLED=0 go1.8 build -a -installsuffix cgo

