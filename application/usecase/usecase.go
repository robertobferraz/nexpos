package usecase

import (
	"context"
	"fmt"
	"github.com/robertobff/food-service/adapter/outbound/auth"
	"github.com/robertobff/food-service/application/dto"
	dtoDomain "github.com/robertobff/food-service/domain/dto"
	"github.com/robertobff/food-service/domain/entity"
	"github.com/robertobff/food-service/domain/repository"
	"github.com/robertobff/food-service/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

var Module = fx.Module(
	"usecase",
	fx.Provide(NewUsecase),
)

type Usecase struct {
	userRepo        repository.UserRepository
	userOrdersRepo  repository.UserOrdersRepository
	userAddressRepo repository.UserAddressRepository
	itemRepo        repository.ItemRepository
	categoryRepo    repository.CategoryRepository
	discountRepo    repository.DiscountRepository
	countryRepo     repository.CountryRepository
	stateRepo       repository.StateRepository
	cityRepo        repository.CityRepository
	districtRepo    repository.DistrictRepository
	streetRepo      repository.StreetRepository
	fb              *auth.AuthFirebase
	logger          *zap.SugaredLogger
}

func NewUsecase(
	userRepo repository.UserRepository,
	userOrdersRepo repository.UserOrdersRepository,
	itemRepo repository.ItemRepository,
	categoryRepo repository.CategoryRepository,
	countryRepo repository.CountryRepository,
	stateRepo repository.StateRepository,
	cityRepo repository.CityRepository,
	districtRepo repository.DistrictRepository,
	streetRepo repository.StreetRepository,
	fb *auth.AuthFirebase,
	logger *zap.SugaredLogger,
) (*Usecase, error) {
	return &Usecase{
		userRepo:       userRepo,
		userOrdersRepo: userOrdersRepo,
		itemRepo:       itemRepo,
		categoryRepo:   categoryRepo,
		countryRepo:    countryRepo,
		stateRepo:      stateRepo,
		cityRepo:       cityRepo,
		districtRepo:   districtRepo,
		streetRepo:     streetRepo,
		fb:             fb,
		logger:         logger,
	}, nil
}

func (u *Usecase) CreateUser(ctx context.Context, idto *dto.CreateUserInDto) (*dto.CreateUserOutDto, error) {
	password, err := utils.CryptPassword(idto.Password)
	if err != nil {
		u.logger.Errorw("error while creating user", "error: ", err)
		return nil, err
	}

	birthDate, err := time.Parse("2006-01-02", *idto.Birthdate)
	if err != nil {
		u.logger.Errorw("error while creating user", "error: ", err)
		return nil, err
	}

	user := &entity.User{
		Email:       idto.Email,
		Username:    idto.Username,
		Name:        idto.Name,
		Password:    idto.Password,
		PhoneNumber: idto.PhoneNumber,
		BirthDate:   utils.PTime(birthDate),
		Cpf:         idto.Cpf,
		ExternalID:  nil,
	}

	existingUser, err := u.fb.GetUserByEmail(ctx, user.Email)
	if err != nil {
		u.logger.Error("Error verifying email: ", err)
		return nil, err
	}
	if existingUser != nil {
		u.logger.Info("Attempted registration with existing email address: ", *user.Email)
		return nil, fmt.Errorf("email %s is already in use", *user.Email)
	}

	userFire, err := u.fb.CreateUser(ctx, user)
	if err != nil {
		u.logger.Errorw("Error creating user no firebase", "error", err)
		return nil, err
	}

	user, err = entity.NewUser(user.Name, user.Username, user.Email, password, user.Cpf, idto.PhoneNumber, idto.Birthdate, userFire.ExternalID)
	if err != nil {
		u.logger.Errorw("error while creating user", "error: ", err)
		return nil, err
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		u.logger.Errorw("error while creating user", "error: ", err)
		return nil, err
	}

	response := &dto.CreateUserOutDto{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return response, nil
}

func (u *Usecase) DeleteUser(ctx context.Context, idto *dto.DeleteUserInDto) error {
	user, err := u.userRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     idto.ID,
			},
		},
	})
	if err != nil {
		u.logger.Errorw("error while deleting user", "error: ", err)
		return err
	}

	if user == nil {
		u.logger.Warnw("user not found", "id", idto.ID)
		return nil
	}

	err = u.userRepo.Delete(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     user.ID,
			},
		},
	})

	if err != nil {
		u.logger.Errorw("error while deleting user", "error: ", err)
		return err
	}

	return nil
}

func (u *Usecase) CreateCategory(ctx context.Context, idto *dto.CreateCategoryInDto) (*dto.CreateCategoryOutDto, error) {
	category, err := entity.NewCategory(idto.Name, idto.Description, idto.Image)
	if err != nil {
		u.logger.Errorw("error while creating category", "error: ", err)
		return nil, err
	}

	err = u.categoryRepo.Create(ctx, category)
	if err != nil {
		u.logger.Errorw("error while creating category", "error: ", err)
		return nil, err
	}

	response := &dto.CreateCategoryOutDto{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Image:       category.Image,
	}

	return response, nil
}

func (u *Usecase) DeleteCategory(ctx context.Context, idto *dto.DeleteCategoryInDto) error {
	category, err := u.categoryRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     idto.ID,
			},
		},
	})
	if err != nil {
		u.logger.Errorw("error while deleting category", "error: ", err)
		return err
	}

	if category == nil {
		u.logger.Warnw("category not found", "id", idto.ID)
		return nil
	}

	err = u.categoryRepo.Delete(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     category.ID,
			},
		},
	})

	if err != nil {
		u.logger.Errorw("error while deleting category", "error: ", err)
		return err
	}

	return nil
}

func (u *Usecase) CreateItem(ctx context.Context, idto *dto.CreateItemInDto) (*dto.CreateItemOutDto, error) {
	category, err := u.categoryRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     idto.CategoryID,
			},
		},
	})

	if err != nil {
		u.logger.Errorw("error while creating item", "error: ", err)
		return nil, err
	}

	if category == nil {
		u.logger.Warnw("category not found", "id", idto.CategoryID)
		return nil, nil
	}

	item, err := entity.NewItem(idto.Name, idto.Description, idto.Image, idto.Price, category)
	if err != nil {
		u.logger.Errorw("error while creating item", "error: ", err)
		return nil, err
	}

	err = u.itemRepo.Create(ctx, item)
	if err != nil {
		u.logger.Errorw("error while creating item", "error: ", err)
		return nil, err
	}

	categoryDto := &dto.CreateCategoryOutDto{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Image:       category.Image,
	}

	response := &dto.CreateItemOutDto{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
		Category:    categoryDto,
	}

	return response, nil
}

func (u *Usecase) DeleteItem(ctx context.Context, idto *dto.DeleteItemInDto) error {
	item, err := u.itemRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     idto.ID,
			},
		},
	})

	if err != nil {
		u.logger.Errorw("error while deleting item", "error: ", err)
		return err
	}

	if item == nil {
		u.logger.Warnw("item not found", "id", idto.ID)
		return nil
	}

	err = u.itemRepo.Delete(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     idto.ID,
			},
		},
	})

	if err != nil {
		u.logger.Errorw("error while deleting item", "error: ", err)
		return err
	}

	return nil
}

func (u *Usecase) CreateDiscount(ctx context.Context, idto *dto.CreateDiscountInDto) (*dto.CreateDiscountOutDto, error) {
	category := &entity.Category{}
	item := &entity.Item{}
	var err error
	if idto.ItemID != nil {
		item, err = u.itemRepo.Find(ctx, &dtoDomain.GormQuery{
			Where: &[]dtoDomain.GormWhere{
				{
					Column:    "id",
					Condition: "=",
					Value:     idto.ItemID,
				},
			},
		})
		if err != nil {
			u.logger.Errorw("error while creating discount", "error: ", err)
			return nil, err
		}
	}

	if idto.CategoryID != nil {
		category, err = u.categoryRepo.Find(ctx, &dtoDomain.GormQuery{
			Where: &[]dtoDomain.GormWhere{
				{
					Column:    "id",
					Condition: "=",
					Value:     idto.CategoryID,
				},
			},
		})
		if err != nil {
			u.logger.Errorw("error while creating discount", "error: ", err)
			return nil, err
		}
	}

	discount, err := entity.NewDiscount(category, item, idto.Date, idto.Value)
	if err != nil {
		u.logger.Errorw("error while creating discount", "error: ", err)
		return nil, err
	}

	itemDto := &dto.CreateItemOutDto{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Image:       item.Image,
	}

	categoryDto := &dto.CreateCategoryOutDto{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Image:       category.Image,
	}

	response := &dto.CreateDiscountOutDto{
		ID:       discount.ID,
		Item:     itemDto,
		Category: categoryDto,
		Value:    discount.Value,
		Date:     discount.Date,
	}

	return response, nil
}

func (u *Usecase) DeleteDiscount(ctx context.Context, idto *dto.DeleteDiscountInDto) error {
	discount, err := u.discountRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     idto.ID,
			},
		},
	})

	if err != nil {
		u.logger.Errorw("error while deleting discount", "error: ", err)
		return err
	}

	if discount == nil {
		u.logger.Warnw("discount not found", "id", idto.ID)
		return nil
	}

	err = u.discountRepo.Delete(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     idto.ID,
			},
		},
	})

	if err != nil {
		u.logger.Errorw("error while deleting discount", "error: ", err)
		return err
	}

	return nil
}
