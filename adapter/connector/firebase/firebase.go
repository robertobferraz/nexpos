package firebase

import (
	"context"
	"encoding/json"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

var Module = fx.Module(
	"auth",
	ConfigModule,
	fx.Provide(NewFirebase),
)

type Firebase struct {
	app    *firebase.App
	client *auth.Client
	ctx    context.Context
	logger *zap.SugaredLogger
}

func NewFirebase(c *Config, logger *zap.SugaredLogger) (*Firebase, error) {
	ctx := context.Background()

	var opts []option.ClientOption
	if *c.Env == "local" {
		base, mErr := json.Marshal(c)
		if mErr != nil {
			logger.Fatal(mErr)
		}

		opts = append(opts, option.WithCredentialsJSON(base))
	}

	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: *c.ProjectId,
	}, opts...)
	if err != nil {
		logger.Errorf("failed to initialize Firebase app: %v", err)
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		logger.Errorf("failed to initialize Firebase auth client: %v", err)
		return nil, err
	}

	return &Firebase{
		app:    app,
		client: client,
		ctx:    ctx,
		logger: logger,
	}, nil
}

func (f *Firebase) Client() *auth.Client {
	return f.client
}
