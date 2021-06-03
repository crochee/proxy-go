package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetBalanceNode godoc
// @Summary GetBalanceNode
// @Description 获取负载节点
// @Security ApiKeyAuth
// @Tags Balance
// @Accept application/json
// @Produce  application/json
// @Success 200 {string} string "ok"
// @Failure 400 {object} e.Response
// @Failure 500 {object} e.Response
// @Router /v1/nodes [get]
func GetBalanceNode(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ok")
}
