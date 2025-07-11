package conversion

import "elevate-hub/types/dto"

func Department2Child(departments []*dto.Department, root *dto.Department) *dto.Department {
	// 使用map存储部门，提高查找效率
	departmentMap := make(map[int64]*dto.Department)
	for _, dept := range departments {
		departmentMap[dept.ID] = dept
	}
	if root == nil {
		root = &dto.Department{
			ID: 0,
		}
	}

	for _, dept := range departments {
		if dept.Parent == root.ID {
			root.Children = append(root.Children, dept)
			dept.Children = []*dto.Department{}
			buildDepartmentChildren(departmentMap, dept)
		}
	}
	return root
}

func buildDepartmentChildren(departmentMap map[int64]*dto.Department, parent *dto.Department) {
	for _, dept := range departmentMap {
		if dept.Parent == parent.ID {
			parent.Children = append(parent.Children, dept)
			dept.Children = []*dto.Department{}
			buildDepartmentChildren(departmentMap, dept)
		}
	}
}
