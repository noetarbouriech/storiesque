package img

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	S3 *minio.Client
}

func NewMinio(host, accessKeyID, secretAccessKey string) (*Minio, error) {

	// initialize minio client
	minioClient, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	// create bucket
	err = minioClient.MakeBucket(context.Background(), "storiesque", minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(context.Background(), "storiesque")
		if errBucketExists == nil && exists {
			log.Println("Bucket storiesque already exists")
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Println("Successfully created bucket storiesque")
	}

	// set public policy on minio
	policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::storiesque/*"],"Sid": ""}]}`
	err = minioClient.SetBucketPolicy(context.Background(), "storiesque", policy)
	if err != nil {
		log.Fatalln(err)
	}

	return &Minio{S3: minioClient}, nil
}
