run-metadata:
	@go run metadata/cmd/*.go
run-rating:
	@go run rating/cmd/main.go
run-movie:
	@go run movie/cmd/*.go
protoc:
	@protoc -I=api --go_out=. --go-grpc_out=. movie.proto
benchmark:
	@go test cmd/sizecompare/*.go -bench=.
