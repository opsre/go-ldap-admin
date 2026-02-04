package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/eryajf/go-ldap-admin/controller"
	"github.com/gin-gonic/gin"
)

// LoginHandler
// @Summary 登录接口 (手动加上: Bearer + token(密码加密接口))
// @Description 用户登录
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param  data body request.RegisterAndLoginReq true "用户登录信息账号和密码"
// @Success 200 {object} response.ResponseBody
// @Router /base/login [post]
func LoginHandler() {}

// LogoutHandler
// @Summary 退出登录
// @Description 用户退出登录
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseBody
// @Router /base/logout [post]
func LogoutHandler() {
}

// RefreshHandler
// @Summary 刷新 Token
// @Description 使用旧的 Token 获取新的 Token
// @Tags 基础管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 旧的 Token"
// @Success 200 {object} response.ResponseBody
// @Router /base/refreshToken [post]
func RefreshHandler() {

}

// 注册基础路由
func InitBaseRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	base := r.Group("/base")
	{
		base.GET("ping", controller.Demo)
		base.GET("encryptpwd", controller.Base.EncryptPasswd) // 生成加密密码
		base.GET("decryptpwd", controller.Base.DecryptPasswd) // 密码解密为明文
		base.GET("config", controller.Base.GetConfig)         // 获取系统配置
		// 登录登出刷新token无需鉴权
		base.POST("/login", authMiddleware.LoginHandler)
		base.POST("/logout", authMiddleware.LogoutHandler)
		base.POST("/refreshToken", authMiddleware.RefreshHandler)
		base.POST("/sendcode", controller.Base.SendCode)   // 给用户邮箱发送验证码
		base.POST("/changePwd", controller.Base.ChangePwd) // 修改用户密码
		base.GET("/dashboard", controller.Base.Dashboard)  // 系统首页展示数据
	}
	return r
}
