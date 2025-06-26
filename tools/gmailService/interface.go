package gmailService

import (
	"context"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

type ClientInterface interface {
	Authenticate(ctx context.Context, clientSecretPath string, port int) (*oauth2.Token, error)
	CreateGmailService(ctx context.Context, credentialsPath, tokenPath string) (*gmail.Service, error)
}
