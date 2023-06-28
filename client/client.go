package client

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/scalescape/dolores/config"
	"github.com/scalescape/dolores/store/google"
)

type Client struct {
	Service
	bucket string
	ctx    context.Context //nolint:containedctx
}

type EncryptedConfig struct {
	Environment string `json:"environment"`
	Name        string `json:"name"`
	Data        string `json:"data"`
}

func (c *Client) UploadSecrets(req EncryptedConfig) error {
	log.Trace().Msgf("uploading to %s name: %s", c.bucket, req.Name)
	return c.Service.Upload(c.ctx, req, c.bucket)
}

type FetchSecretRequest struct {
	Environment string `json:"environment"`
	Name        string `json:"name"`
}
type FetchSecretResponse struct {
	Data string `json:"data"`
}

func (c *Client) FetchSecrets(req FetchSecretRequest) ([]byte, error) {
	data, err := c.Service.FetchConfig(c.ctx, c.bucket, req)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type Recipient struct {
	PublicKey string `json:"public_key"`
}

type OrgPublicKeys struct {
	Recipients []Recipient `json:"recipients"`
}

func (c *Client) GetOrgPublicKeys(env string) (OrgPublicKeys, error) {
	keys, err := c.Service.GetOrgPublicKeys(c.ctx, env, c.bucket)
	if err != nil || len(keys) == 0 {
		return OrgPublicKeys{}, err
	}
	recps := make([]Recipient, len(keys))
	for i, k := range keys {
		recps[i].PublicKey = k
	}
	return OrgPublicKeys{Recipients: recps}, nil
}

func New(ctx context.Context, cfg config.Client) (*Client, error) {
	if err := cfg.Valid(); err != nil {
		return nil, err
	}
	gcfg := google.Config{ServiceAccountFile: cfg.Google.ApplicationCredentials}
	st, err := google.NewStore(ctx, gcfg)
	if err != nil {
		return nil, err
	}
	return &Client{ctx: ctx, Service: Service{store: st}, bucket: cfg.BucketName()}, nil
}