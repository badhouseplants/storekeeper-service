package uploader

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	proto_uploader "github.com/droplez/droplez-go-proto/pkg/uploader"
	"github.com/droplez/droplez-uploader/tools/constants"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implementation of the uploader grpc service
type grpcUploaderImpl struct {
	proto_uploader.UnimplementedUploaderServer
	Unlocked chan bool
}

// Register grpc service
func Register(grpcServer *grpc.Server, unlocked chan bool) {
	proto_uploader.RegisterUploaderServer(grpcServer, &grpcUploaderImpl{
		Unlocked: unlocked,
	})
}

/*
* Upload file to Minio
* This function stores the file locally and then put it in the Minio bucket
* First we need to check that system has enough space to store the file
* * If yes -> just store a file, put ut in the bucket, and remove
* * If no -> put the request in the queue and wait until the system has enough free space
* Requests should respect the order and not be proccessed out of turn
 */
func (s grpcUploaderImpl) Upload(stream proto_uploader.Uploader_UploadServer) (err error) {
	var (
		fileName             = uuid.New().String() // temporary file name for storing locally
		minioEndpoint        = viper.GetString(constants.MinioEndpoint)
		minioAccessKeyID     = viper.GetString(constants.MinioAccessKeyID)
		minioSecretAccessKey = viper.GetString(constants.MinioSecretAccessKey)
		minioBucket          = viper.GetString(constants.MinioBucket)
		minioUseSSL          = false
	)

	// Receive metadata first
	chunk, err := stream.Recv()
	if err != nil {
		return
	} else if chunk.Content != nil {
		return status.Error(codes.InvalidArgument, constants.ErrOnlyMetadataAllowed)
	}

	meta := chunk.GetFileMetadata()
	var (
		// If name is empty, generate a new name for minio object
		minioFileName = uuid.New().String() + filepath.Ext(meta.GetLocalName())
		minioFileSize = meta.GetFileSize()
		minioUserID   = meta.GetUserId()
	)

	// Initialize minio client object
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKeyID, minioSecretAccessKey, ""),
		Secure: minioUseSSL,
	})
	if err != nil {
		return err
	}

	// Create a temporary local file to save payload
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	// If file is created, it need to be removed on error or on return
	defer cleanUp(f)

	/*
	* This variable is used to write file iteratively
	* It will store the point from which the next chunk will be written
	 */
	offset := int64(0)

	// Receive a file
	for {
		chunk, err := stream.Recv()

		// Quit the loop if the whole file is received
		if err == io.EOF {
			break
		}

		// Write chunk to the file
		byte := chunk.GetContent()
		offsetTemp, err := f.WriteAt(byte, offset)
		if err != nil {
			return err
		}
		offset += int64(offsetTemp)

		// If file is bigger than the claimed filesize return error
		if offset > minioFileSize {
			return status.Error(codes.InvalidArgument, constants.ErrWrongFileSizeProviced)
		}

	}

	// Put the file in the Minio bucket
	_, err = minioClient.FPutObject(stream.Context(), minioBucket, minioFileName, fileName, minio.PutObjectOptions{
		PartSize: uint64(minioFileSize),
		UserTags: map[string]string{
			"user-id": minioUserID,
		},
	})
	if err != nil {
		log.Print(err)
		return err
	}

	return stream.SendAndClose(&proto_uploader.UploadedFileData{
		Object: minioFileName,
	})
}

func (s grpcUploaderImpl) GetDownloadableLink(ctx context.Context, in *proto_uploader.UploadedFileData) (*proto_uploader.DownloadableLink, error) {
	var (
		minioEndpoint        = viper.GetString(constants.MinioEndpoint)
		minioAccessKeyID     = viper.GetString(constants.MinioAccessKeyID)
		minioSecretAccessKey = viper.GetString(constants.MinioSecretAccessKey)
		minioBucket          = viper.GetString(constants.MinioBucket)
		minioUseSSL          = false
	)
	fmt.Println(in)
	// Initialize minio client object
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKeyID, minioSecretAccessKey, ""),
		Secure: minioUseSSL,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Get downloadable link from minio
	// info, err := minioClient.PresignedGetObject(ctx, minioBucket, in.Object, time.Minute, url.Values{})
	info, err := minioClient.PresignedGetObject(ctx, minioBucket, in.Object, time.Minute, nil)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &proto_uploader.DownloadableLink{
		Url: info.String(),
	}, nil

}

/*
* Helpers
 */

// Remove file from local storage after uploading to Minio or on error
func cleanUp(file *os.File) {
	if err := file.Close(); err != nil {
		log.Print(err)
	}
	if err := os.Remove(file.Name()); err != nil {
		log.Print(err)
	}
}
