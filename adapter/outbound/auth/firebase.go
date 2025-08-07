package auth

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/robertobff/nexpos/application/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	firebase "firebase.google.com/go/v4/auth"
	fb "github.com/robertobff/nexpos/adapter/connector/firebase"
	"github.com/robertobff/nexpos/domain/entity"
	fb_errors "github.com/robertobff/nexpos/domain/errors"
	"github.com/robertobff/nexpos/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"firebase",
	fx.Provide(NewAuthFirebase),
)

type Firebase struct {
	firebaseClient *fb.Firebase
	logger         *zap.SugaredLogger
}

type DeletionRequest struct {
	UID         *string    `firestore:"uid"`
	CreatedAt   *time.Time `firestore:"created_at"`
	ExpiresAt   *time.Time `firestore:"expires_at"`
	CronJobID   *int       `firestore:"cron_job_id"`
	IsCompleted *bool      `firestore:"is_completed"`
	Cancelled   *bool      `firestore:"cancelled"`
}

func NewAuthFirebase(firebaseClient *fb.Firebase, logger *zap.SugaredLogger) *Firebase {
	return &Firebase{
		firebaseClient: firebaseClient,
		logger:         logger,
	}
}

func (a *Firebase) SetFirestore(document *string, v DeletionRequest) error {
	_, err := a.firebaseClient.Firestore().Collection("deletion_request").Doc(*document).Set(context.Background(), v)
	return err
}

func (a *Firebase) UpdateFirestore(document *string, v []firestore.Update) error {
	_, err := a.firebaseClient.Firestore().Collection("deletion_request").Doc(*document).Update(context.Background(), v)
	return err
}

func (a *Firebase) GetDeletionRequest(ctx context.Context, id *string) (*DeletionRequest, error) {
	doc, err := a.firebaseClient.Firestore().Collection("deletion_requests").Doc(*id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, err
	}

	var req DeletionRequest
	if err := doc.DataTo(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (a *Firebase) Iter(ctx context.Context) *firestore.DocumentIterator {
	iter := a.firebaseClient.Firestore().Collection("deletion_requests").Documents(ctx)
	return iter
}

func (a *Firebase) UpdateUserPassword(ctx context.Context, uid string, password string) error {
	params := (&firebase.UserToUpdate{}).Password(password)

	_, err := a.firebaseClient.Client().UpdateUser(ctx, uid, params)
	if err != nil {
		return err
	}

	return nil
}

func (a *Firebase) VerifyIDToken(ctx context.Context, idToken *string) (*firebase.Token, error) {
	return a.firebaseClient.Client().VerifyIDToken(ctx, *idToken)
}

func (a *Firebase) VerifyAccessToken(ctx context.Context, t *string) (*string, error) {
	token, err := a.firebaseClient.Client().VerifyIDToken(ctx, *t)
	if err != nil {
		return nil, fb_errors.ErrInvalidToken
	}
	uid := token.UID

	return utils.PString(uid), nil
}

func (a *Firebase) FindUserByUID(ctx context.Context, uid *string) (*entity.User, error) {
	var user *entity.User

	u, err := a.firebaseClient.Client().GetUser(ctx, *uid)
	if err != nil {
		return user, err
	}

	user.Email = utils.PString(u.Email)
	user.ExternalID = utils.PString(u.UID)
	user.Username = utils.PString(u.DisplayName)

	//user.Password = utils.PString(fmt.Sprintf("randpass:%s", uuid.New()))

	return user, nil
}

func (a *Firebase) GetUserByEmail(ctx context.Context, email *string) (*entity.User, error) {
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

func (a *Firebase) CreateUser(ctx context.Context, user *dto.CreateUserInDto) (*dto.CreateUserInDto, error) {
	params := (&firebase.UserToCreate{}).
		Email(*user.Email).
		DisplayName(*user.Name).
		PhoneNumber(*user.PhoneNumber).
		Password(*user.Password).
		Disabled(false)

	u, err := a.firebaseClient.Client().CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	user.ExternalID = utils.PString(u.UID)

	return user, nil
}

func (a *Firebase) DisableUser(ctx context.Context, uid *string) error {
	params := (&firebase.UserToUpdate{}).Disabled(true)

	_, err := a.firebaseClient.Client().UpdateUser(ctx, *uid, params)
	if err != nil {
		a.logger.Error("error disabling user:", err)
		return err
	}

	a.logger.Info("user with UID: ", *uid, ", successfully disabled")
	return nil
}

func (a *Firebase) EnableUser(ctx context.Context, uid *string) error {
	params := (&firebase.UserToUpdate{}).Disabled(false)

	_, err := a.firebaseClient.Client().UpdateUser(ctx, *uid, params)
	if err != nil {
		a.logger.Error("error re-enabling user: ", err)
		return err
	}

	a.logger.Info("User with UID: ", *uid, ", successfully re-enabled")
	return nil
}

func (a *Firebase) IsUserEnabled(ctx context.Context, uid *string) (*bool, error) {
	u, err := a.firebaseClient.Client().GetUser(ctx, *uid)
	if err != nil {
		if firebase.IsUserNotFound(err) {
			a.logger.Info("user with UID: ", *uid, ",  not found")
			return utils.PBool(false), fmt.Errorf("user with UID %s not found", *uid)
		}
		a.logger.Error("error searching for user: ", err)
		return utils.PBool(false), err
	}

	isEnabled := !u.Disabled
	a.logger.Info("user with UID: ", *uid, ", is enabled: ", isEnabled)
	return utils.PBool(isEnabled), nil
}
