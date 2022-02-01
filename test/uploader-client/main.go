package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	proto_uploader "github.com/badhouseplants/storekeeper-go-proto/pkg"
	"google.golang.org/grpc"
)

func main() {
	// var opt []grpc.DialOption
	conn, err := grpc.Dial("0.0.0.0:9090", grpc.WithInsecure())
	if err != nil {
		return
	}
	// defer conn.Close()
	client := proto_uploader.NewUploaderClient(conn)

	fmt.Println("INIT GRPC")
	stream, err := client.Upload(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	
		fmt.Println("UPLOADING")
	arch, err := os.Open("./test/uploader-client/test.zip")
	if err != nil {
		fmt.Println(err)
	}

	fileSize, err := arch.Stat()
	if err != nil {
		fmt.Println(err)
	}

	// First sending metadata without content
	stream.Send(&proto_uploader.Chunk{
		Content: nil,
		FileMetadata: &proto_uploader.Metadata{
			// Name:        conf.Project.GetObjectName(),
			ContentType: proto_uploader.Metadata_CONTENT_TYPE_ARCHIVE,
			LocalName:   "./resources/testfile.zip",
			FileSize:    fileSize.Size(),
			UserId:      "user id",
		},
	})

	r := bufio.NewReader(arch)
	for {
		bytes, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		err = stream.Send(&proto_uploader.Chunk{
			Content: bytes,
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	reply, err := stream.CloseAndRecv()
	fmt.Println(reply.GetObject())

	// data, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	}
}
