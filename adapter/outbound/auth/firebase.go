package auth

import (
	"context"
	firebase "firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/google/uuid"
	fb "github.com/robertobff/food-service/adapter/connector/firebase"
	"github.com/robertobff/food-service/domain/entity"
	firebase_errors "github.com/robertobff/food-service/domain/errors"
	"github.com/robertobff/food-service/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"firebase",
	fx.Provide(NewAuthFirebase),
)

type AuthFirebase struct {
	firebaseClient *fb.Firebase
	logger         *zap.SugaredLogger
}

func NewAuthFirebase(firebaseClient *fb.Firebase, logger *zap.SugaredLogger) *AuthFirebase {
	return &AuthFirebase{
		firebaseClient: firebaseClient,
		logger:         logger,
	}
}

func (a *AuthFirebase) UpdateUserPassword(ctx context.Context, uid string, password string) error {
	params := (&firebase.UserToUpdate{}).Password(password)

	_, err := a.firebaseClient.Client().UpdateUser(ctx, uid, params)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthFirebase) VerifyIDToken(ctx context.Context, idToken *string) (*firebase.Token, error) {
	return a.firebaseClient.Client().VerifyIDToken(ctx, *idToken)
}

func (a *AuthFirebase) VerifyAccessToken(ctx context.Context, t *string) (*string, error) {
	token, err := a.firebaseClient.Client().VerifyIDToken(ctx, *t)
	if err != nil {
		return nil, firebase_errors.ErrInvalidToken
	}
	uid := token.UID

	return utils.PString(uid), nil
}

func (a *AuthFirebase) FindUserByUID(ctx context.Context, uid *string) (*entity.User, error) {
	var user *entity.User

	u, err := a.firebaseClient.Client().GetUser(ctx, *uid)
	if err != nil {
		return user, err
	}

	user.Email = utils.PString(u.Email)
	user.ExternalID = utils.PString(u.UID)
	user.Username = utils.PString(u.DisplayName)

	user.Password = utils.PString(fmt.Sprintf("randpass:%s", uuid.New()))

	return user, nil
}

func (a *AuthFirebase) GetUserByEmail(ctx context.Context, email *string) (*entity.User, error) {
	var user entity.User

	u, err := a.firebaseClient.Client().GetUserByEmail(ctx, *email)
	if err != nil {
		if firebase.IsUserNotFound(err) {
			a.logger.Info("Email %s not found", *email)
			return nil, nil
		}
		a.logger.Error("Error searching for user by email: %v", err)
		return nil, err
	}

	user.Email = utils.PString(u.Email)
	user.ExternalID = utils.PString(u.UID)
	user.Username = utils.PString(u.DisplayName)

	return &user, nil
}

func (a *AuthFirebase) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	params := (&firebase.UserToCreate{}).
		Email(*user.Email).
		Password(*user.Password).
		DisplayName(*user.Name).
		PhoneNumber(*user.PhoneNumber).
		Disabled(false)

	u, err := a.firebaseClient.Client().CreateUser(ctx, params)
	if err != nil {
		return user, err
	}

	user.ExternalID = utils.PString(u.UID)

	return user, nil
}
