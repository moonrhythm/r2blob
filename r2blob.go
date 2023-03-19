package r2blob

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
)

func init() {
	blob.DefaultURLMux().RegisterBucket(Scheme, new(URLOpener))
}

const Scheme = "r2"

type URLOpener struct {
}

func (o *URLOpener) OpenBucketURL(ctx context.Context, u *url.URL) (*blob.Bucket, error) {
	bucketName := u.Host
	account := u.Query().Get("account")
	accessKeyID := u.Query().Get("access_key_id")
	accessKeySecret := u.Query().Get("access_key_secret")

	if bucketName == "" {
		return nil, fmt.Errorf("r2: missing bucket name")
	}
	if account == "" {
		return nil, fmt.Errorf("r2: missing account name")
	}
	if accessKeyID == "" {
		return nil, fmt.Errorf("r2: missing access key id")
	}
	if accessKeySecret == "" {
		return nil, fmt.Errorf("r2: missing access key secret")
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", account),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return s3blob.OpenBucketV2(ctx, client, bucketName, nil)
}
