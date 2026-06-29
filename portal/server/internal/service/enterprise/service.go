package enterprise

import (
	"portal/internal/config"
	"portal/internal/repository"
)

type Service struct {
	Cache              repository.Cache
	Config             config.Config
	UserRepo           repository.UserRepo
	EnterpriseRepo     repository.EnterpriseRepo
	UserEnterpriseRepo repository.UserEnterpriseRepo
}
