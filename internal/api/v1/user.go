package v1

import (
	"gin-mini-starter/internal/middleware"
	"gin-mini-starter/internal/model"
	"gin-mini-starter/pkg/resp"

	"github.com/gin-gonic/gin"
)

type LoginResp struct {
	Token string `json:"token" example:"xxx"`
}

type LoginUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 用户登陆 godoc
// @Summary 用户登陆
// @Description 用户登陆
// @Tags 用户
// @Accept json
// @Produce json
// @Param User body LoginUserReq true "用户登陆信息"
// @Success 200 {object} resp.Resp{data=LoginResp}
// @Router /User/login [post]
func LoginUser(c *gin.Context) {
	var req LoginUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, resp.CodeParamsInvalid, err.Error())
		return
	}
	User := model.User{}
	err := User.Query().Where("username = ?", req.Username).First(&User).Error
	if err != nil {
		resp.Error(c, resp.CodeInternalServer, "用户名或密码错误")
		return
	}
	if !User.ComparePassword(req.Password) {
		resp.Error(c, resp.CodeInternalServer, "用户名或密码错误")
		return
	}
	// 去除敏感信息
	User.Password = ""
	// 生成JWT
	token, err := middleware.GenerateToken(middleware.UserClaims{
		ID:       User.ID,
		UserName: User.Username,
		Role:     []middleware.Role{middleware.RoleUser},
		UserType: "User",
	})
	if err != nil {
		resp.Error(c, resp.CodeInternalServer, err.Error())
		return
	}
	resp.SuccessData(c, gin.H{
		"token": token,
	})
}

type UpdateUserPasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// 用户修改密码 godoc
// @Summary 用户修改密码
// @Description 用户修改密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param req body UpdateUserPasswordReq true "用户修改密码信息"
// @Success 200 {object} resp.Resp
// @Router /User/password [post]
func UpdateUserPassword(c *gin.Context) {
	var req UpdateUserPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error(c, resp.CodeParamsInvalid, err.Error())
		return
	}
	// 获取用户ID
	user := c.MustGet("user").(*middleware.UserClaims)
	User := model.User{}
	err := User.Query().Where("id = ?", user.ID).First(&User).Error
	if err != nil {
		resp.Error(c, resp.CodeInternalServer, err.Error())
		return
	}
	if !User.ComparePassword(req.OldPassword) {
		resp.Error(c, resp.CodeInternalServer, "旧密码错误")
		return
	}
	User.Password = req.NewPassword
	User.EncryptPassword()
	err = User.Query().Updates(&User).Error
	if err != nil {
		resp.Error(c, resp.CodeInternalServer, err.Error())
		return
	}
	resp.Success(c)
}
