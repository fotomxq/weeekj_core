package Router2Mid

import (
	"github.com/gin-gonic/gin"
)

//url封装体系

// RouterGlob 全局通用结构方法存储
type RouterGlob struct {
	//路由结构体
	Routers *gin.RouterGroup
	//路由级别
	level string
}

// RouterURL 方法覆盖
type RouterURL struct {
	//全局方法
	BaseData RouterGlob
}

// Top 根URL
func (t *RouterURL) Top(urlPath string) RouterURL {
	//映射URL
	return RouterURL{
		BaseData: RouterGlob{
			Routers: t.BaseData.Routers.Group(urlPath, headerBaseData),
			level:   "top",
		},
	}
}

// Base 基础封装
func (t *RouterURL) Base(urlPath string) *RouterURL {
	//映射URL
	return &RouterURL{
		BaseData: RouterGlob{
			Routers: t.BaseData.Routers.Group(urlPath),
			level:   "base",
		},
	}
}

// RouterURLPublic 有头级别路由
type RouterURLPublic struct {
	//全局方法
	BaseData RouterGlob
}

// Public 无头封装
func (t *RouterURL) Public() *RouterURLPublic {
	//映射URL
	return &RouterURLPublic{
		BaseData: RouterGlob{
			Routers: t.BaseData.Routers.Group("/public"),
			level:   "public",
		},
	}
}

// RouterURLIOT IOT级别路由
type RouterURLIOT struct {
	//全局方法
	BaseData RouterGlob
}

// IOT 无头封装
func (t *RouterURL) IOT() *RouterURLIOT {
	//映射URL
	return &RouterURLIOT{
		BaseData: RouterGlob{
			Routers: t.BaseData.Routers.Group("/iot"),
			level:   "iot",
		},
	}
}

// RouterURLHeader 有头级别路由
type RouterURLHeader struct {
	//全局方法
	BaseData RouterGlob
}

// Header 有头封装
func (t *RouterURL) Header() *RouterURLHeader {
	//映射URL
	return &RouterURLHeader{
		BaseData: RouterGlob{
			Routers: t.BaseData.Routers.Group("/header", headerLoginBefore),
			level:   "header",
		},
	}
}

// RouterURLUser 用户级别路由
type RouterURLUser struct {
	//全局方法
	BaseData RouterGlob
}

// User 用户级别封装
func (t *RouterURL) User() *RouterURLUser {
	//映射URL
	return &RouterURLUser{
		BaseData: RouterGlob{
			Routers: t.BaseData.Routers.Group("/user", headerLoggedUser),
			level:   "user",
		},
	}
}

// RouterURLManager 管理层级路由
type RouterURLManager struct {
	//全局方法
	BaseData RouterGlob
}

// Manager 管理级别封装
func (t *RouterURL) Manager() *RouterURLUser {
	//映射URL
	return &RouterURLUser{
		BaseData: RouterGlob{
			Routers: t.BaseData.Routers.Group("/manager", headerLoggedUser),
			level:   "manager",
		},
	}
}

// RouterURLRole 用户角色级别路由
type RouterURLRole struct {
	//全局方法
	BaseData RouterGlob
	//角色类型
	RoleType string
}

// Role 用户级别封装
func (t *RouterURL) Role(roleType string) *RouterURLRole {
	//映射URL
	return &RouterURLRole{
		BaseData: RouterGlob{
			Routers: t.BaseData.Routers.Group("/role", headerLoggedUser),
			level:   "role",
		},
		RoleType: roleType,
	}
}

// RouterURLOrg 组织级别路由
type RouterURLOrg struct {
	//全局方法
	BaseData RouterGlob
}

// OB 组织成员级别封装
func (t *RouterURL) OB() *RouterURLOrg {
	//映射URL
	return &RouterURLOrg{
		BaseData: RouterGlob{
			Routers: t.BaseData.Routers.Group("/ob", headerLoggedUser, headerSelectedOrg),
			level:   "ob",
		},
	}
}

// OM 组织管理级别封装
func (t *RouterURL) OM() *RouterURLOrg {
	//映射URL
	return &RouterURLOrg{
		BaseData: RouterGlob{
			Routers: t.BaseData.Routers.Group("/om", headerLoggedUser, headerSelectedOrg),
			level:   "om",
		},
	}
}
