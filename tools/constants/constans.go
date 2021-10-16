package constants

// Viper keys
const (
	// Minio
	MinioBucket          = "MINIO_BUCKET"
	MinioEndpoint        = "MINIO_ENDPOINT"
	MinioAccessKeyID     = "MINIO_ACCESS_KEY_ID"
	MinioSecretAccessKey = "MINIO_SERCET_ACCESS_KEY"
	// Droplez
	DroplezUploaderPort = "DROPLEZ_UPLOADER_PORT"
	DroplezUploaderMode = "DROPLEZ_UPLOADER_MODE"
)

// Errors
const (
	ErrInternalMinioError    = "internal minio error"
	ErrWrongFileSizeProviced = "wrong file size provided"
	ErrOnlyMetadataAllowed   = "only metadate is allowed in the first message"
)