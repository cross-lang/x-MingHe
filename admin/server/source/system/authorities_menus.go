package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderMenuAuthority = initOrderMenu + initOrderAuthority

type initMenuAuthority struct{}

// auto run
func init() {
	system.RegisterInit(initOrderMenuAuthority, &initMenuAuthority{})
}

func (i *initMenuAuthority) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil // do nothing
}

func (i *initMenuAuthority) TableCreated(ctx context.Context) bool {
	return false // always replace
}

func (i *initMenuAuthority) InitializerName() string {
	return "sys_menu_authorities"
}

func (i *initMenuAuthority) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	initAuth := &initAuthority{}
	authorities, ok := ctx.Value(initAuth.InitializerName()).([]sysModel.SysAuthority)
	if !ok {
		return ctx, errors.Wrap(system.ErrMissingDependentContext, "创建 [菜单-权限] 关联失败, 未找到权限表初始化数据")
	}

	allMenus, ok := ctx.Value(new(initMenu).InitializerName()).([]sysModel.SysBaseMenu)
	if !ok {
		return next, errors.Wrap(errors.New(""), "创建 [菜单-权限] 关联失败, 未找到菜单表初始化数据")
	}
	next = ctx

	// 构建菜单ID映射，方便快速查找
	menuMap := make(map[uint]sysModel.SysBaseMenu)
	for _, menu := range allMenus {
		menuMap[menu.ID] = menu
	}

	// 为不同角色分配不同权限
	// 1. 超级管理员角色(888) - 拥有所有菜单权限
	if err = db.Model(&authorities[0]).Association("SysBaseMenus").Replace(allMenus); err != nil {
		return next, errors.Wrap(err, "为超级管理员分配菜单失败")
	}

	// 2. 开发者(666) - 仅拥系统工具、插件系统、示例文件、服务器状态、关于我们、超级管理员（角色管理、菜单管理、api管理、字典管理、参数管理）
	// 仅选择部分父级菜单及其子菜单
	menu666Set := map[string]struct{}{
		"person":                            {}, // 个人信息
		"systemTools":                       {}, // 系统工具
		"autoPkg":                           {}, // 模板管理
		"autoCode":                          {}, // 代码生成器
		"autoCodeAdmin":                     {}, // 自动化代码管理
		"system":                            {}, // 系统配置
		"exportTemplate":                    {}, // 导出模板
		"sysVersion":                        {}, // 版本管理
		"sysError":                          {}, // 错误日志
		"plugin":                            {}, // 插件系统
		"https://plugin.gin-vue-admin.com/": {}, // 插件市场
		"installPlugin":                     {}, // 插件系统
		"pubPlug":                           {}, // 打包插件
		"plugin-email":                      {}, // 邮件插件
		"anInfo":                            {}, // 公告管理[示例]
		"example":                           {}, // 示例文件
		"upload":                            {}, // 媒体库（上传下载）
		"breakpoint":                        {}, // 断点续传
		"customer":                          {}, // 客户列表（资源示例）
		"state":                             {}, // 服务器状态
		"about":                             {}, // 关于我们
		"superAdmin":                        {}, // 超级管理员
		"authority":                         {}, // 角色管理
		"menu":                              {}, // 菜单管理
		"api":                               {}, // api管理
		"dictionary":                        {}, // 字典管理
		"sysParams":                         {}, // 参数管理
	}
	var menu666 []sysModel.SysBaseMenu
	for _, menu := range allMenus {
		if _, ok := menu666Set[menu.Name]; ok {
			menu666 = append(menu666, menu)
		}
	}

	if err = db.Model(&authorities[1]).Association("SysBaseMenus").Replace(menu666); err != nil {
		return next, errors.Wrap(err, "为普通用户分配菜单失败")
	}
	return next, nil
}

func (i *initMenuAuthority) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	auth := &sysModel.SysAuthority{}
	if ret := db.Model(auth).
		Where("authority_id = ?", 666).Preload("SysBaseMenus").Find(auth); ret != nil {
		if ret.Error != nil {
			return false
		}
		return len(auth.SysBaseMenus) > 0
	}
	return false
}
