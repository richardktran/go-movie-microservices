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

# Build binary file
build-metadata:
	GOOS=linux go build -o metadata/main metadata/cmd/*.go

build-rating:
	GOOS=linux go build -o rating/main rating/cmd/*.go

build-movie:
	GOOS=linux go build -o movie/main movie/cmd/*.go

# Build docker images
docker-build-metadata:
	@docker build -t metadata -f metadata/Dockerfile ./metadata
docker-build-rating:
	@docker build -t rating -f rating/Dockerfile ./rating
docker-build-movie:
	@docker build -t movie -f movie/Dockerfile ./movie

# Run docker images
docker-run-metadata:
	@docker run -p 8081:8081 -it metadata
docker-run-rating:
	@docker run -p 8082:8082 -it rating
docker-run-movie:
	@docker run -p 8083:8083 -it movie