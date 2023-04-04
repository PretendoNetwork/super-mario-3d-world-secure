package globals

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type PresignClient struct {
	cfg aws.Config
}

type PostObjectInput struct {
	// Key name
	Key string

	// The name of the bucket to presign the post to
	Bucket string

	// Expiration -  The number of seconds the presigned post is valid for.
	ExpiresIn time.Duration

	// A list of conditions to include in the policy. Each element can be either a list or a structure.
	// For example:
	// [
	//      {"acl": "public-read"}, ["content-length-range", 2, 5], ["starts-with", "$success_action_redirect", ""]
	// ]
	Conditions []interface{}
}

type PresignedPostObject struct {
	URL           string `json:"url"`
	Key           string `json:"key"`
	Policy        string `json:"policy"`
	Credential    string `json:"credential"`
	SecurityToken string `json:"securityToken,omitempty"`
	Signature     string `json:"signature"`
	Date          string `json:"date"`
}

func (presignClient *PresignClient) PresignPostObject(input *PostObjectInput) (*PresignedPostObject, error) {
	credentials, err := presignClient.cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		return nil, err
	}

	expirationTime := time.Now().Add(input.ExpiresIn).UTC()
	dateString := expirationTime.Format("20060102")
	dateTimeString := expirationTime.Format("20060102T150405Z")

	credentialString := fmt.Sprintf("%s/%s/%s/s3/aws4_request", credentials.AccessKeyID, dateString, presignClient.cfg.Region)

	policyDocument, err := createPolicyDocument(expirationTime, input.Bucket, input.Key, credentialString, &credentials.SessionToken, input.Conditions)
	if err != nil {
		return nil, err
	}

	signature := createSignature(credentials.SecretAccessKey, presignClient.cfg.Region, dateString, policyDocument)

	// * This is the best way, that I can find, to format the
	// * s3 endpoint URL without needing to pass it in twice
	// * (first duirng s3 setup, then when creating this client)
	exp := regexp.MustCompile(`^(https://)(.*)`)
	endpoint, err := presignClient.cfg.EndpointResolverWithOptions.ResolveEndpoint("", "")
	if err != nil {
		return nil, err
	}

	url := exp.ReplaceAllString(endpoint.URL, fmt.Sprintf("${1}%s.${2}", input.Bucket))

	presignedPostObject := &PresignedPostObject{
		Key:           input.Key,
		Policy:        policyDocument,
		Signature:     signature,
		URL:           url,
		Credential:    credentialString,
		SecurityToken: credentials.SessionToken,
		Date:          dateTimeString,
	}

	return presignedPostObject, nil
}

func NewPresignClient(cfg aws.Config) *PresignClient {
	return &PresignClient{cfg}
}

func createPolicyDocument(expirationTime time.Time, bucket string, key string, credentialString string, securityToken *string, extraConditions []interface{}) (string, error) {

	doc := map[string]interface{}{}
	doc["expiration"] = expirationTime.Format("2006-01-02T15:04:05.000Z")

	// conditions
	conditions := []interface{}{}
	conditions = append(conditions, map[string]string{
		"bucket": bucket,
	})

	conditions = append(conditions, []string{
		"starts-with", "$key", key,
	})

	conditions = append(conditions, map[string]string{
		"x-amz-credential": credentialString,
	})

	if securityToken != nil {
		conditions = append(conditions, map[string]string{
			"x-amz-security-token": *securityToken,
		})
	}

	conditions = append(conditions, map[string]string{
		"x-amz-algorithm": "AWS4-HMAC-SHA256",
	})

	conditions = append(conditions, map[string]string{
		"x-amz-date": expirationTime.Format("20060102T150405Z"),
	})

	// other conditions
	conditions = append(conditions, extraConditions...)

	doc["conditions"] = conditions

	// base64 encoded json string
	jsonBytes, err := json.Marshal(doc)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

func createSignature(secretKey string, region string, dateString string, stringToSign string) string {

	// Helper to make the HMAC-SHA256.
	makeHmac := func(key []byte, data []byte) []byte {
		hash := hmac.New(sha256.New, key)
		hash.Write(data)
		return hash.Sum(nil)
	}

	h1 := makeHmac([]byte("AWS4"+secretKey), []byte(dateString))
	h2 := makeHmac(h1, []byte(region))
	h3 := makeHmac(h2, []byte("s3"))
	h4 := makeHmac(h3, []byte("aws4_request"))
	signature := makeHmac(h4, []byte(stringToSign))
	return hex.EncodeToString(signature)
}
