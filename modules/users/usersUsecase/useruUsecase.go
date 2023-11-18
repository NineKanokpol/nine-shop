package usersUsecase

import (
	"fmt"

	"github.com/NineKanokpol/Nine-shop-test/config"
	"github.com/NineKanokpol/Nine-shop-test/modules/users"
	usersRepositories "github.com/NineKanokpol/Nine-shop-test/modules/users/usersRepositories"
	"github.com/NineKanokpol/Nine-shop-test/pkg/nineauth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUseCase interface {
	InsertCustomer(req *users.UserRegisterRequest) (*users.UserPassport, error)
	GetPassport(req *users.UserCredentials) (*users.UserPassport, error)
	RefreshTokenPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUserRespository
}

func UsersUseCase(cfg config.IConfig, usersRepository usersRepositories.IUserRespository) IUsersUseCase {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: usersRepository,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterRequest) (*users.UserPassport, error) {
	//hashing password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	result, err := u.usersRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *usersUsecase) GetPassport(req *users.UserCredentials) (*users.UserPassport, error) {
	//Find user
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	//Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}

	//signToken
	accessToken, err := nineauth.NewNineAuth(nineauth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	refreshToken, err := nineauth.NewNineAuth(nineauth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	//set passport
	passport := &users.UserPassport{
		User: &users.User{
			ID:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}
	if err := u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}
	return passport, nil
}

func (u *usersUsecase) RefreshTokenPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	//Parse Token
	claims, err := nineauth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	//check oauth
	oauth, err := u.usersRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	//Find profile
	profile, err := u.usersRepository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	//*sign payload
	newClaims := &users.UserClaims{
		Id:     profile.ID,
		RoleId: profile.RoleId,
	}

	accessToken, err := nineauth.NewNineAuth(
		nineauth.Access,
		u.cfg.Jwt(),
		newClaims,
	)
	if err != nil {
		return nil, err
	}

	refreshToken := nineauth.RepeatToken(
		u.cfg.Jwt(),
		newClaims,
		claims.ExpiresAt.Unix(),
	)

	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}
	if err := u.usersRepository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}
	return passport, nil
}
