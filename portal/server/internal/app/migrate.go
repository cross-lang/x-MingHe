package app

import "portal/internal/model"

func (a *App) AutoMigrate() {
	if !a.Config.Mysql.AutoMigrate {
		return
	}
	a.Container.MysqlDB.AutoMigrate(
		model.XEnterprise{},
		model.XUser{},
		model.XUserEnterprise{},
		model.XUserIdentityVerification{},
	)
}
