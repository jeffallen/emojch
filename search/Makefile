VER=3

docker: search
	docker build -t jeffallen/emojch-search:$(VER) .
	docker push jeffallen/emojch-search:$(VER)

search: *.go
	CGO_ENABLED=0 go1.8 build -a -installsuffix cgo

