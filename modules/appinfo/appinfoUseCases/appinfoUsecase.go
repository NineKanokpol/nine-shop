package appinfoUseCases

import (
	"github.com/NineKanokpol/Nine-shop-test/modules/appinfo"
	"github.com/NineKanokpol/Nine-shop-test/modules/appinfo/appinfoRepositories"
)

type IAppinfoUsecase interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
}

type appinfoUsecase struct {
	appinfoRepository appinfoRepositories.IAppinfoRepository
}

func AppinfoUseCase(appinfoRepository appinfoRepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{appinfoRepository: appinfoRepository}
}

func (u *appinfoUsecase) FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {

	category, err := u.appinfoRepository.FindCategory(req)
	if err != nil {
		return nil, err
	}

	return category, nil
}
