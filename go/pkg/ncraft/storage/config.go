package storage

// Config describes an object storage backend. MinIO also supports other
// S3-compatible servers through the "s3" vendor alias.
type Config struct {
	Vendor       string `json:"vendor" default:"minio"`
	Endpoint     string `json:"endpoint"`
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	BucketName   string `json:"bucketName"`
	Region       string `json:"region"`
	SessionToken string `json:"sessionToken"`
	// Secure controls TLS for bare host:port endpoints (default false for compatibility).
	// An explicit http:// or https:// endpoint takes precedence.
	Secure *bool `json:"secure"`
	// Prefix namespaces all object keys within the bucket.
	Prefix string `json:"prefix"`
	// CreateBucket defaults to true for compatibility. Set false for an existing
	// bucket with credentials restricted to object reads and writes.
	CreateBucket *bool `json:"createBucket"`
}
