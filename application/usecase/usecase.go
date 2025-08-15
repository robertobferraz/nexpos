package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/robertobff/nexpos/adapter/outbound/api/openCep"
	"github.com/robertobff/nexpos/adapter/outbound/auth"
	"github.com/robertobff/nexpos/adapter/outbound/scheduler"
	"github.com/robertobff/nexpos/application/dto"
	dtoDomain "github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
	"github.com/robertobff/nexpos/domain/repository"
	"github.com/robertobff/nexpos/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
	imageRepo       repository.ImageRepository
	openCep         *openCep.Api
	fb              *auth.Firebase
	schedule        *scheduler.Scheduler
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
	imageRepo repository.ImageRepository,
	userAddressRepo repository.UserAddressRepository,
	openCep *openCep.Api,
	fb *auth.Firebase,
	schedule *scheduler.Scheduler,
	logger *zap.SugaredLogger,
) (*Usecase, error) {
	return &Usecase{
		userRepo:        userRepo,
		userOrdersRepo:  userOrdersRepo,
		itemRepo:        itemRepo,
		categoryRepo:    categoryRepo,
		countryRepo:     countryRepo,
		stateRepo:       stateRepo,
		cityRepo:        cityRepo,
		districtRepo:    districtRepo,
		streetRepo:      streetRepo,
		imageRepo:       imageRepo,
		userAddressRepo: userAddressRepo,
		openCep:         openCep,
		fb:              fb,
		schedule:        schedule,
		logger:          logger,
	}, nil
}

func (u *Usecase) CreateUserIfNotExist(ctx context.Context, idto *dto.CreateUserInDto) (*entity.User, error) {
	user, err := entity.NewUser(idto.Name, idto.Username, idto.Email, idto.Cpf, idto.PhoneNumber, nil, idto.ExternalID, nil)
	if err != nil {
		u.logger.Errorw("error while creating user", "error: ", err)
		return nil, err
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		u.logger.Errorw("error while creating user", "error: ", err)
		return nil, err
	}

	return user, nil
}

func (u *Usecase) CheckUserDeletion(ctx context.Context, user *entity.User) error {
	err := u.schedule.CancelUserDeletion(ctx, user.ID)
	if err != nil {
		u.logger.Errorw("error while canceling user", "error: ", err)
		return err
	}

	return nil
}

func (u *Usecase) CreateUser(ctx context.Context, idto *dto.CreateUserInDto) (*entity.User, error) {
	existingUser, err := u.fb.GetUserByEmail(ctx, idto.Email)
	if err != nil {
		u.logger.Error("error verifying email: ", err)
		return nil, err
	}
	if existingUser != nil {
		u.logger.Info("attempted registration with existing email address: ", *idto.Email)
		return nil, fmt.Errorf("email %s is already in use", *idto.Email)
	}

	var birthDate time.Time
	if idto.Birthdate != nil {
		birthDate, err = time.Parse("2006-01-02", *idto.Birthdate)
		if err != nil {
			u.logger.Errorw("error while creating user", "error: ", err)
			return nil, err
		}
	}

	userFire, err := u.fb.CreateUser(ctx, idto)
	if err != nil {
		u.logger.Errorw("error creating user on firebase", "error", err)
		return nil, err
	}

	user, err := entity.NewUser(idto.Name, idto.Username, idto.Email, idto.Cpf, idto.PhoneNumber, utils.PString(birthDate.Format("2006-01-02")), userFire.ExternalID, nil)
	if err != nil {
		u.logger.Errorw("error while creating user", "error: ", err)
		return nil, err
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		u.logger.Errorw("error while creating user", "error: ", err)
		return nil, err
	}

	return user, nil
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

	if user.DeletedAt != nil {
		u.logger.Warnw("user is deleted", "id", idto.ID)
		return nil
	}

	user.DeletedAt = utils.PTime(time.Now())
	err = u.userRepo.Save(ctx, user)
	if err != nil {
		u.logger.Errorw("error while deleting user", "error: ", err)
		return err
	}

	err = u.schedule.ScheduleUserDeletion(ctx, user.ID, user.ExternalID)
	if err != nil {
		u.logger.Errorw("error while deleting user", "error: ", err)
		return err
	}

	return nil
}

func (u *Usecase) SaveUser(ctx context.Context, idto *dto.SaveUserInDto) (*entity.User, error) {
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
		u.logger.Errorw("error while saving user", "error: ", err)
		return nil, err
	}

	if user == nil {
		u.logger.Warnw("user not found", "id", idto.ID)
		return nil, errors.New("user not found")
	}

	if idto.Username != nil {
		if err := user.SetUsername(idto.Username); err != nil {
			u.logger.Errorw("error while saving user", "error: ", err)
			return nil, err
		}
	}

	if idto.Name != nil {
		if err := user.SetName(idto.Name); err != nil {
			u.logger.Errorw("error while saving user", "error: ", err)
			return nil, err
		}
	}

	if idto.Birthdate != nil {
		if err := user.SetBirthDate(idto.Birthdate); err != nil {
			u.logger.Errorw("error while saving user", "error: ", err)
			return nil, err
		}
	}

	if idto.PhoneNumber != nil {
		if err := user.SetPhoneNumber(idto.PhoneNumber); err != nil {
			u.logger.Errorw("error while saving user", "error: ", err)
			return nil, err
		}
	}

	if idto.Image != nil {
		img, err := entity.NewImage(idto.Image.Name, idto.Image.ContentType, idto.Image.Data)
		if err != nil {
			u.logger.Errorw("error while saving user", "error: ", err)
			return nil, err
		}

		if err := user.SetImage(img); err != nil {
			u.logger.Errorw("error while saving user", "error: ", err)
			return nil, err
		}
	}

	if idto.Email != nil {
		if err := user.SetEmail(idto.Email); err != nil {
			u.logger.Errorw("error while saving user", "error: ", err)
			return nil, err
		}
	}

	err = u.userRepo.Save(ctx, user)
	if err != nil {
		u.logger.Errorw("error while saving user", "error: ", err)
		return nil, err
	}

	return user, nil
}

func (u *Usecase) SaveUserInAuth(ctx context.Context, user *entity.User) error {
	err := u.userRepo.Save(ctx, user)
	if err != nil {
		u.logger.Errorw("error while saving user", "error: ", err)
		return err
	}

	return nil
}

func (u *Usecase) GetImage(ctx context.Context, id *string) ([]byte, *string, error) {
	image, err := u.imageRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     id,
			},
		},
	})

	if err != nil {
		return nil, nil, err
	}

	if image == nil {
		return nil, nil, errors.New("user not found")
	}

	return *image.Data, image.ContentType, nil
}

func (u *Usecase) GetUserByUID(ctx context.Context, idto *dto.GetUserByUIDInDto) (*entity.User, error) {
	user, err := u.userRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "external_id",
				Condition: "=",
				Value:     idto.UID,
			},
		},
		Unscoped: true,
	})

	if err != nil {
		u.logger.Errorw("error while getting user", "error: ", err)
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return user, nil
}

func (u *Usecase) GetUsers(ctx context.Context, idto *dto.GetUsersInDto) (*[]dto.GetUsersOutDto, error) {
	users, err := u.userRepo.Get(ctx, &dtoDomain.GormQuery{
		Preload: &[]dtoDomain.GormPreload{{
			Field: "Image",
		}},
	})
	if err != nil {
		u.logger.Errorw("error while getting users", "error: ", err)
		return nil, err
	}

	var resp []dto.GetUsersOutDto
	for _, value := range *users {
		var user dto.GetUsersOutDto
		if value.Image != nil {
			user.Image = &dto.Image{
				Name:        value.Image.Name,
				Url:         utils.PString(fmt.Sprintf("%s://%s/v1/image/%s", *idto.Protocol, *idto.HostName, *value.ImageID)),
				ContentType: value.Image.ContentType,
			}
		}
		user.ID = value.ID
		user.Name = value.Name
		user.Email = value.Email
		user.Birthdate = value.BirthDate
		user.PhoneNumber = value.PhoneNumber
		user.Username = value.Username

		resp = append(resp, user)
	}

	return &resp, nil
}

func (u *Usecase) CreateCategory(ctx context.Context, idto *dto.CreateCategoryInDto) (*entity.Category, error) {
	var image *entity.Image
	var err error
	if idto.Image != nil {
		image, err = entity.NewImage(idto.Image.Name, idto.Image.ContentType, idto.Image.Data)
		if err != nil {
			u.logger.Errorw("error while creating category", "error: ", err)
			return nil, err
		}
	}

	category, err := entity.NewCategory(idto.Name, idto.Description, image)
	if err != nil {
		u.logger.Errorw("error while creating category", "error: ", err)
		return nil, err
	}

	err = u.categoryRepo.Create(ctx, category)
	if err != nil {
		u.logger.Errorw("error while creating category", "error: ", err)
		return nil, err
	}

	return category, nil
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

func (u *Usecase) CreateItem(ctx context.Context, idto *dto.CreateItemInDto) (*entity.Item, error) {
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

	var image *entity.Image
	if idto.Image != nil {
		image, err = entity.NewImage(idto.Image.Name, idto.Image.ContentType, idto.Image.Data)
		if err != nil {
			u.logger.Errorw("error while creating category", "error: ", err)
			return nil, err
		}
	}

	item, err := entity.NewItem(idto.Name, idto.Description, idto.Price, category, image)
	if err != nil {
		u.logger.Errorw("error while creating item", "error: ", err)
		return nil, err
	}

	err = u.itemRepo.Create(ctx, item)
	if err != nil {
		u.logger.Errorw("error while creating item", "error: ", err)
		return nil, err
	}

	return item, nil
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

func (u *Usecase) CreateDiscount(ctx context.Context, idto *dto.CreateDiscountInDto) (*entity.Discount, error) {
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

	return discount, nil
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

func (u *Usecase) CreateCountry(ctx context.Context, idto *dto.CreateCountryInDto) error {
	country, err := entity.NewCountry(idto.Name, idto.Iso2, idto.Iso3, idto.PhoneCode, idto.Capital, idto.CurrencyCode, idto.Emoji, idto.ExternalID)
	if err != nil {
		u.logger.Errorw("error while creating country", "error: ", err)
		return err
	}

	err = u.countryRepo.Create(ctx, country)
	if err != nil {
		u.logger.Errorw("error while creating country", "error: ", err)
		return err
	}

	return nil
}

func (u *Usecase) CreateState(ctx context.Context, idto *dto.CreateStateInDto) error {
	country, err := u.countryRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     idto.CountryID,
			},
		},
	})
	if err != nil {
		u.logger.Errorw("error while creating state", "error: ", err)
		return err
	}

	if country == nil {
		u.logger.Warnw("country not found", "id", idto.CountryID)
		return errors.New("country not found")
	}

	state, err := entity.NewState(idto.Name, idto.Iso2, idto.ExternalID, country)
	if err != nil {
		u.logger.Errorw("error while creating state", "error: ", err)
		return err
	}

	err = u.stateRepo.Create(ctx, state)
	if err != nil {
		u.logger.Errorw("error while creating state", "error: ", err)
		return err
	}

	return nil
}

func (u *Usecase) CreateCity(ctx context.Context, idto *dto.CreateCityInDto) error {
	state, err := u.stateRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     idto.StateID,
			},
		},
	})

	if err != nil {
		u.logger.Errorw("error while creating city", "error: ", err)
		return err
	}

	if state == nil {
		u.logger.Warnw("state not found", "id", idto.StateID)
		return errors.New("state not found")
	}

	city, err := entity.NewCity(idto.Name, idto.ExternalID, state)
	if err != nil {
		u.logger.Errorw("error while creating city", "error: ", err)
		return err
	}

	err = u.cityRepo.Create(ctx, city)
	if err != nil {
		u.logger.Errorw("error while creating city", "error: ", err)
		return err
	}

	return nil
}

func (u *Usecase) CreteDistrict(ctx context.Context, idto *dto.CreateDistrictInDto) error {
	city, err := u.cityRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     idto.CityID,
			},
		},
	})

	if err != nil {
		u.logger.Errorw("error while creating district", "error: ", err)
		return err
	}

	if city == nil {
		u.logger.Warnw("city not found", "id", idto.CityID)
		return errors.New("city not found")
	}

	district, err := entity.NewDistrict(idto.Name, city)
	if err != nil {
		u.logger.Errorw("error while creating district", "error: ", err)
		return err
	}

	err = u.districtRepo.Create(ctx, district)
	if err != nil {
		u.logger.Errorw("error while creating district", "error: ", err)
		return err
	}

	return nil
}

func (u *Usecase) CreateStreet(ctx context.Context, idto *dto.CreateStreetInDto) error {
	district, err := u.districtRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     idto.DistrictID,
			},
		},
	})

	if err != nil {
		u.logger.Errorw("error while creating street", "error: ", err)
		return err
	}

	if district == nil {
		u.logger.Warnw("district not found", "id", idto.DistrictID)
		return errors.New("district not found")
	}

	street, err := entity.NewStreet(idto.Name, idto.ZipCode, idto.Number, district)
	if err != nil {
		u.logger.Errorw("error while creating street", "error: ", err)
		return err
	}

	err = u.streetRepo.Create(ctx, street)
	if err != nil {
		u.logger.Errorw("error while creating street", "error: ", err)
		return err
	}

	return nil
}

func (u *Usecase) FindCountryByIdentifier(ctx context.Context, idto *dto.FindCountryByIdentifierInDto) (*entity.Country, error) {
	country, err := u.countryRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     idto.Identifier,
			},
		},
	})

	if err != nil {
		u.logger.Errorw("error while finding country", "error: ", err)
		return nil, err
	}

	if country == nil {
		u.logger.Warnw("country not found", "id", idto.Identifier)
		return nil, errors.New("country not found")
	}

	return country, nil
}

func (u *Usecase) CreateUserAddressCondition(ctx context.Context, idto *dto.CreateUserAddressInDto) (*entity.UserAddress, error) {
	country, err := u.countryRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     idto.CountryID,
			},
		},
	})
	if err != nil {
		u.logger.Errorw("error while creating user address", "error", err)
		return nil, err
	}

	if country == nil {
		u.logger.Warnw("country not found", "id", idto.CountryID)
		return nil, errors.New("country not found")
	}

	switch *country.Iso2 {
	case "BR":
		if idto.ActionBrazil == nil {
			u.logger.Error("missing payload for Brazil")
			return nil, errors.New("missing payload for Brazil")
		}
		return u.actionBrazil(ctx, idto.ActionBrazil)
	default:
		if idto.ActionOthers == nil {
			u.logger.Error("missing payload for other countries")
			return nil, errors.New("missing payload for other countries")
		}
		return u.actionOthers(ctx, idto.ActionOthers)
	}
}

func (u *Usecase) actionBrazil(ctx context.Context, idto *dto.ActionBrazilInDto) (*entity.UserAddress, error) {
	country, err := u.countryRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     idto.CountryID,
			},
		},
	})

	if err != nil {
		u.logger.Errorw("error while finding country", "error: ", err)
		return nil, err
	}

	user, err := u.userRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "id",
				Condition: "=",
				Value:     idto.UserID,
			},
		},
	})

	if err != nil {
		u.logger.Errorw("error while creating district", "error: ", err)
		return nil, err
	}

	if user == nil {
		u.logger.Warnw("user not found", "id", idto.UserID)
		return nil, errors.New("user not found")
	}

	resp, err := u.openCep.GetByCep(idto.Cep)
	if err != nil {
		u.logger.Errorw("error while getting brazil", "error: ", err)
		return nil, err
	}

	if resp != nil {
		state, err := u.stateRepo.Find(ctx, &dtoDomain.GormQuery{
			Where: &[]dtoDomain.GormWhere{
				{
					Column:    "country_id",
					Condition: "=",
					Value:     country.ID,
				},
				{
					Column:    "iso2",
					Condition: "=",
					Value:     resp.Uf,
				},
			},
		})
		if err != nil {
			u.logger.Errorw("error while getting brazil", "error: ", err)
			return nil, err
		}
		state.Country = country
		city, err := u.cityRepo.Find(ctx, &dtoDomain.GormQuery{
			Where: &[]dtoDomain.GormWhere{
				{
					Column:    "state_id",
					Condition: "=",
					Value:     state.ID,
				},
				{
					Column:    "unaccent(name)",
					Condition: "ILIKE unaccent(?)",
					Value:     utils.PString("%" + *resp.Localidade + "%"),
				},
			},
		})

		if err != nil {
			u.logger.Errorw("error while getting brazil", "error: ", err)
			return nil, err
		}

		if city == nil {
			u.logger.Warnw("city not found", "id", idto.UserID)
			return nil, errors.New("city not found")
		}

		district, err := u.districtRepo.Find(ctx, &dtoDomain.GormQuery{
			Where: &[]dtoDomain.GormWhere{
				{
					Column:    "city_id",
					Condition: "=",
					Value:     city.ID,
				},
				{
					Column:    "unaccent(name)",
					Condition: "ILIKE unaccent(?)",
					Value:     utils.PString("%" + *resp.Bairro + "%"),
				},
			},
			Debug: true,
		})
		if err != nil {
			u.logger.Errorw("error while getting brazil", "error: ", err)
			return nil, err
		}

		if district == nil {
			newDistrict, err := entity.NewDistrict(resp.Bairro, city)
			if err != nil {
				u.logger.Errorw("error while creating district", "error: ", err)
				return nil, err
			}

			err = u.districtRepo.Create(ctx, newDistrict)
			if err != nil {
				u.logger.Errorw("error while creating district", "error: ", err)
				return nil, err
			}

			district = newDistrict
		}

		street, err := u.streetRepo.Find(ctx, &dtoDomain.GormQuery{
			Where: &[]dtoDomain.GormWhere{
				{
					Column:    "district_id",
					Condition: "=",
					Value:     district.ID,
				},
				{
					Column:    "unaccent(name)",
					Condition: "ILIKE unaccent(?)",
					Value:     utils.PString("%" + *resp.Logradouro + "%"),
				},
			},
			Preload: &[]dtoDomain.GormPreload{
				{
					Field: "District.City.State.Country",
				},
			},
		})

		if err != nil {
			u.logger.Errorw("error while getting brazil", "error: ", err)
			return nil, err
		}

		existStreet := true
		var newStreet *entity.Street
		if street == nil {
			existStreet = false
			newStreet, err = entity.NewStreet(resp.Logradouro, idto.Cep, idto.Number, district)
			if err != nil {
				u.logger.Errorw("error while creating street", "error: ", err)
				return nil, err
			}
		} else if street.Number != idto.Number {
			newStreet, err = entity.NewStreet(resp.Logradouro, idto.Cep, idto.Number, district)
			if err != nil {
				u.logger.Errorw("error while creating street", "error: ", err)
				return nil, err
			}
		} else if street.ZipCode != idto.Cep {
			newStreet, err = entity.NewStreet(resp.Logradouro, idto.Cep, idto.Number, district)
			if err != nil {
				u.logger.Errorw("error while creating street", "error: ", err)
				return nil, err
			}
		}

		if !existStreet {
			err = u.streetRepo.Create(ctx, newStreet)
			if err != nil {
				u.logger.Errorw("error while creating street", "error: ", err)
				return nil, err
			}

			street = newStreet
			street.District.City = city
			street.District.City.State = state
		}

		newUserAddress, err := entity.NewUserAddress(user, street)
		if err != nil {
			u.logger.Errorw("error while creating user address", "error: ", err)
			return nil, err
		}

		err = u.userAddressRepo.Create(ctx, newUserAddress)
		if err != nil {
			u.logger.Errorw("error while creating user address", "error: ", err)
			return nil, err
		}

		return newUserAddress, nil
	} else {
		return nil, errors.New("cep not found")
	}
}

func (u *Usecase) actionOthers(ctx context.Context, idto *dto.ActionOthersInDto) (*entity.UserAddress, error) {
	return nil, nil
}

func (u *Usecase) GetUserAddress(ctx context.Context, idto *dto.GetUserAddressInDto) (*dto.GetUserAddressOutDto, error) {
	userAddress, err := u.userAddressRepo.Find(ctx, &dtoDomain.GormQuery{
		Where: &[]dtoDomain.GormWhere{
			{
				Column:    "user_id",
				Condition: "=",
				Value:     idto.UserID,
			},
		},
		Preload: &[]dtoDomain.GormPreload{
			{
				Field: "Street.District.City.State.Country",
			},
			{
				Field: "User",
			},
		},
	})

	if err != nil {
		u.logger.Errorw("error while getting user address", "error: ", err)
		return nil, err
	}

	if userAddress == nil {
		u.logger.Warnw("user address not found", "id", idto.UserID)
		return nil, errors.New("user address not found")
	}

	var response dto.GetUserAddressOutDto
	var image *dto.Image
	if userAddress.User.Image != nil {
		image = &dto.Image{
			Name:        userAddress.User.Image.Name,
			Url:         utils.PString(fmt.Sprintf("%s://%s/v1/image/%s", *idto.Protocol, *idto.HostName, *userAddress.User.ImageID)),
			ContentType: userAddress.User.Image.ContentType,
		}
	}

	user := dto.GetUsersOutDto{
		ID:          userAddress.UserID,
		Username:    userAddress.User.Username,
		Name:        userAddress.User.Name,
		Email:       userAddress.User.Email,
		Birthdate:   userAddress.User.BirthDate,
		PhoneNumber: userAddress.User.PhoneNumber,
		Image:       image,
	}

	response.User = &user
	response.Street = userAddress.Street

	return &response, nil
}
