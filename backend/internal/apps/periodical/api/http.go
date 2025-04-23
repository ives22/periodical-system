package api

import (
	"fmt"
	"log"
	"periodical/internal/apps/periodical"
	"periodical/internal/apps/token"
	"periodical/internal/middleware"
	"periodical/internal/response"

	"github.com/gin-gonic/gin"
)

type PeriodicalApiHandler struct {
	// 依赖控制器
	svc periodical.Service
}

func NewPeriodicalApiHandler(svc periodical.Service) *PeriodicalApiHandler {
	return &PeriodicalApiHandler{
		svc: svc,
	}
}

func (h *PeriodicalApiHandler) Registry(r gin.IRouter, tokenSvc token.Service) {
	v1 := r.Group("v1").Group("periodical")
	v1.Use(middleware.Auth(tokenSvc)) // 添加认证中间件
	// POST /v1/periodical
	v1.POST("", h.CreatePeriodicalHandler)
	// POST /v1/periodical/list
	v1.POST("/list", h.QueryPeriodicalHandler)
	// GET /v1/periodical/4
	v1.GET("/:id", h.DescribePeriodicalHandler)
	// PUT /v1/periodical/4
	v1.PUT("/:id", h.UpdatePeriodicalHandler)
	// DELETE /v1/periodical/4
	v1.DELETE("/:id", h.DeletePeriodicalHandler)

	v1.POST("/base/list", h.QueryCategorizeHandler)

}

// CreatePeriodicalHandler 创建期刊
func (h *PeriodicalApiHandler) CreatePeriodicalHandler(ctx *gin.Context) {
	in := periodical.NewCreatePeriodicalRequest()

	if err := ctx.ShouldBindJSON(in); err != nil {
		log.Println(err)
		response.Failed(ctx, response.CodeInvalidParam)
		return
	}

	fmt.Println("============j")
	log.Println(in.BasePeriodical)
	ins, err := h.svc.CreatePeriodical(ctx.Request.Context(), in)
	if err != nil {
		response.Failed(ctx, response.CodeServerBusy)
		return
	}

	response.Success(ctx, ins)
}

// UpdatePeriodicalHandler 更新期刊 （全量）
func (h *PeriodicalApiHandler) UpdatePeriodicalHandler(ctx *gin.Context) {
	log.Println(ctx.Param("id"))
	in := periodical.NewUpdatePeriodicalRequest(ctx.Param("id"))
	log.Println(in.CreatePeriodicalRequest)
	if err := ctx.BindJSON(in.CreatePeriodicalRequest); err != nil {
		log.Println("aaaa")
		log.Println(in.CreatePeriodicalRequest)
		response.Failed(ctx, response.CodeInvalidParam)
		return
	}

	log.Println(in.BasePeriodical)
	ins, err := h.svc.UpdatePeriodical(ctx.Request.Context(), in)
	if err != nil {
		response.Failed(ctx, response.CodeServerBusy)
		return
	}

	response.Success(ctx, ins)
}

// QueryPeriodicalHandler 查询期刊列表
func (h *PeriodicalApiHandler) QueryPeriodicalHandler(ctx *gin.Context) {
	in := periodical.NewQueryPeriodicalRequest()

	if err := ctx.ShouldBind(in); err != nil {
		response.Failed(ctx, response.CodeInvalidParam)
		return
	}

	// 验证排序字段
	if err := in.ValidateOrderBy(); err != nil {
		response.FailedWhitMsg(ctx, response.CodeInvalidParam, err.Error())
		return
	}

	log.Println(in)
	set, err := h.svc.QueryPeriodical(ctx.Request.Context(), in)
	if err != nil {
		response.Failed(ctx, response.CodeServerBusy)
		return
	}

	response.Success(ctx, set)
}

// DescribePeriodicalHandler 查询期刊详情
func (h *PeriodicalApiHandler) DescribePeriodicalHandler(ctx *gin.Context) {
	in := periodical.NewDescribePeriodicalRequest(ctx.Param("id"))

	ins, err := h.svc.DescribePeriodical(ctx.Request.Context(), in)
	if err != nil {
		response.Failed(ctx, response.CodeServerBusy)
		return
	}
	response.Success(ctx, ins)
}

// DeletePeriodicalHandler 删除期刊
func (h *PeriodicalApiHandler) DeletePeriodicalHandler(ctx *gin.Context) {
	in := periodical.NewDeletePeriodicalRequest(ctx.Param("id"))
	err := h.svc.DeletePeriodical(ctx.Request.Context(), in)
	if err != nil {
		response.Failed(ctx, response.CodeServerBusy)
		return
	}
	response.Success(ctx, nil)
}

// QueryCategorizeHandler 查询期刊的各种查询分类列表
func (h *PeriodicalApiHandler) QueryCategorizeHandler(ctx *gin.Context) {
	in := periodical.NewQueryCategorizeRequest()

	if err := ctx.ShouldBind(in); err != nil {
		response.Failed(ctx, response.CodeInvalidParam)
		return
	}

	set, err := h.svc.QueryCategorize(ctx.Request.Context(), in)
	if err != nil {
		response.Failed(ctx, response.CodeServerBusy)
		return
	}

	response.Success(ctx, set)
}
