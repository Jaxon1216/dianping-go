package handler

import (
	"github.com/gin-gonic/gin"
	"go-dianping/api/v1"
	"go-dianping/internal/service"
	"net/http"
)

func NewUserHandler(handler *Handler,
	userService service.UserService,
) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

type UserHandler struct {
	*Handler
	userService service.UserService
}

// SendCode godoc
// @Summary 发送短信验证码并保存验证码
// @Schemes
// @Description
// @Tags user
// @Produce json
// @Param phone query string true "手机号"
// @Success 200 {object} v1.Response
// @Router /user/code [post]
func (h *UserHandler) SendCode(ctx *gin.Context) {
	// var 声明一个变量。这里 req 的类型是 v1.SendCodeReq，
	// 初始值是该 struct 的零值，之后由 ShouldBindQuery 填充。
	var req v1.SendCodeReq

	// 从 URL 的 Query 参数中读取数据，例如：
	// POST /user/code?phone=13800138000
	// &req 表示把 req 的地址传给绑定函数，让函数直接修改 req。
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// Go 常见写法：调用函数、同时声明 err，并立即判断错误。
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.userService.SendCode(ctx.Request.Context(), &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// Login godoc
// @Summary 实现登录功能
// @Schemes
// @Description
// @Tags user
// @Accept json
// @Produce json
// @Param request body v1.LoginReq true "登录请求体"
// @Success 200 {object} v1.LoginResp
// @Router /user/login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var req v1.LoginReq
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	data, err := h.userService.Login(ctx.Request.Context(), &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, data)
}

// Me godoc
// @Summary 获取当前登录的用户并返回
// @Schemes
// @Description
// @Tags user
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.MeResp
// @Router /user/me [get]
func (h *UserHandler) Me(ctx *gin.Context) {
	user, err := h.userService.Me(ctx.Request.Context())
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, user)
}

// QueryUserByID godoc
// @Summary 获取当前登录的用户并返回
// @Schemes
// @Description
// @Tags user
// @Produce json
// @Security Bearer
// @Param id path uint64 true "用户ID"
// @Success 200 {object} v1.MeResp
// @Router /user/me [get]
func (h *UserHandler) QueryUserByID(ctx *gin.Context) {
	var req struct {
		userID uint64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	user, err := h.userService.QueryUserByID(ctx.Request.Context(), req.userID)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, user)
}

// Sign godoc
// @Summary 签到
// @Schemes
// @Description
// @Tags user
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.Response
// @Router /user/sign [post]
func (h *UserHandler) Sign(ctx *gin.Context) {
	err := h.userService.Sign(ctx.Request.Context())
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// SignCount godoc
// @Summary 签到
// @Schemes
// @Description
// @Tags user
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.Response
// @Router /user/sign/count [get]
func (h *UserHandler) SignCount(ctx *gin.Context) {
	err := h.userService.SignCount(ctx.Request.Context())
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
