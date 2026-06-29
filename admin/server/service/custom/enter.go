package custom

import "github.com/flipped-aurora/gin-vue-admin/server/model/system"

type ServiceGroup struct {
	EnterpriseService
	UserService
	system.SysDistrict
}
