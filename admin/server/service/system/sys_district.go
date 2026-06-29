package system

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemResp "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
)

type DistrictService struct {
}

// GetDistrictList 获取区域列表
func (s *DistrictService) GetDistrictList(ctx context.Context, req *systemReq.DistrictList) (*systemResp.DistrictListResponse, error) {
	var districts []system.SysDistrict
	query := global.GVA_DB.WithContext(ctx).Model(&system.SysDistrict{})

	// 父级ID筛选
	if req.PID >= 0 {
		query = query.Where("pid = ?", req.PID)
	}

	// 级别筛选
	if req.Level > 0 {
		query = query.Where("level = ?", req.Level)
	}

	// 关键词搜索
	if req.Keyword != "" {
		query = query.Where("city_name LIKE ?", "%"+req.Keyword+"%")
	}

	// 执行查询
	if err := query.Order("id ASC").Find(&districts).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]systemResp.DistrictDetail, len(districts))
	for i, district := range districts {
		items[i] = systemResp.DistrictDetail{
			ID:       district.ID,
			PID:      district.PID,
			CityName: district.CityName,
			Level:    district.Level,
		}
	}

	return &systemResp.DistrictListResponse{
		Items: items,
	}, nil
}

// GetDistrictTree 获取区域树
func (s *DistrictService) GetDistrictTree(ctx context.Context, req *systemReq.DistrictTree) ([]systemResp.DistrictTreeResponse, error) {
	// 先获取所有区域数据
	var allDistricts []system.SysDistrict
	if err := global.GVA_DB.WithContext(ctx).Find(&allDistricts).Error; err != nil {
		return nil, err
	}

	// 构建区域映射
	districtMap := make(map[int]*system.SysDistrict)
	for i := range allDistricts {
		districtMap[allDistricts[i].ID] = &allDistricts[i]
	}

	// 构建树结构
	tree := buildDistrictTree(req.PID, districtMap)

	return tree, nil
}

// buildDistrictTree 递归构建区域树
func buildDistrictTree(pid int, districtMap map[int]*system.SysDistrict) []systemResp.DistrictTreeResponse {
	var tree []systemResp.DistrictTreeResponse

	// 查找所有子节点
	for id, district := range districtMap {
		if district.PID == pid {
			// 递归构建子树
			children := buildDistrictTree(id, districtMap)

			node := systemResp.DistrictTreeResponse{
				ID:       district.ID,
				PID:      district.PID,
				CityName: district.CityName,
				Level:    district.Level,
			}

			if len(children) > 0 {
				node.Children = children
			}

			tree = append(tree, node)
		}
	}

	return tree
}
