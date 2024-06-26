package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	AuthorizationExpirePeriod = 1800

	APE_ACCESS_KEY = "Access-Key"
	APE_SIGN_TIME  = "Sign-Time"

	// APE time layout
	APE_SIGN_TIME_LAYOUT = "2006-01-02 15:04:05"
)

// Signer abstracts the entity that implements the `Sign` method
type Signer interface {
	// Sign the given Request with the Credentials and SignOptions
	Sign(*http.Request, *ApeCredentials, *SignOptions) (string, error)
}

// SignOptions defines the data structure used by Signer
type SignOptions struct {
	HeadersToSign []string
}

// ApeSigner implements ape sign algorithm
type ApeSigner struct{}

func (a *ApeSigner) Sign(req *http.Request, cred *ApeCredentials, opt *SignOptions) (string, error) {
	// Set the APE request headers
	req.Header.Add(APE_ACCESS_KEY, cred.AccessKey)
	req.Header.Add(APE_SIGN_TIME, util.FormatDate(APE_SIGN_TIME_LAYOUT, util.NowUTCSeconds()))

	// get request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	signInfo := append(body, []byte(cred.SecretAccessKey)...)

	// sign headers
	for _, k := range opt.HeadersToSign {
		if h := req.Header.Get(k); len(h) != 0 {
			signInfo = append(signInfo, []byte(h)...)
		}
	}
	content := sha256.Sum256(signInfo)
	signature := hex.EncodeToString(content[:])
	return signature, nil
}

// 默认使用 DEFAULT_HEADER
var DEFAULT_HEADERS_TO_SIGN = []string{
	strings.ToLower(APE_REQUEST_ID),
	strings.ToLower(APE_SIGN_TIME),
}

var signOptions = &SignOptions{
	HeadersToSign: DEFAULT_HEADERS_TO_SIGN,
}
