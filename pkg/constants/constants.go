package constants

// Viper keys
const (
	// Minio
	ConstMinioBucket          = "MINIO_BUCKET"
	ConstMinioEndpoint        = "MINIO_ENDPOINT"
	ConstMinioAccessKeyID     = "MINIO_ACCESS_KEY_ID"
	ConstMinioSecretAccessKey = "MINIO_SECRET_ACCESS_KEY"
	// Droplez
	ConstDroplezUploaderPort = "DROPLEZ_UPLOADER_PORT"
	ConstDroplezUploaderMode = "DROPLEZ_UPLOADER_MODE"
	ConstDroplezStorePath    = "DROPLEZ_STORE_PATH"
)

// Modes
const (
	ConstDevMode  = "dev"
	ConstProdMode = "prod"
)

// Errors
const (
	ErrInternalMinioError    = "internal minio error"
	ErrWrongFileSizeProviced = "wrong file size provided"
	ErrOnlyMetadataAllowed   = "only metadata is allowed in the first message"
)
