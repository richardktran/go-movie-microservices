package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/richardktran/go-movie-microservices/gen"
	"github.com/richardktran/go-movie-microservices/metadata/pkg/model"
	"google.golang.org/protobuf/proto"
)

var metadata = &model.Metadata{
	ID:          "1",
	Title:       "The Shawshank Redemption",
	Description: "Two imprisoned",
	Director:    "Frank Darabont",
}

var genMetadata = &gen.Metadata{
	Id:          "1",
	Title:       "The Shawshank Redemption",
	Desctiption: "Two imprisoned",
	Director:    "Frank Darabont",
}

func main() {
	jsonBytes, err := serializeToJson(metadata)
	if err != nil {
		panic(err)
	}

	xmlBytes, err := serializeToXml(metadata)
	if err != nil {
		panic(err)
	}

	protoBytes, err := serializeToProto(genMetadata)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JSON size: \t%dB\n", len(jsonBytes))
	fmt.Printf("XML size: \t%dB\n", len(xmlBytes))
	fmt.Printf("Proto size: \t%dB\n", len(protoBytes))
}

func serializeToJson(metadata *model.Metadata) ([]byte, error) {
	return json.Marshal(metadata)
}

func serializeToXml(metadata *model.Metadata) ([]byte, error) {
	return xml.Marshal(metadata)
}

func serializeToProto(metadata *gen.Metadata) ([]byte, error) {
	return proto.Marshal(metadata)
}
