package rpc

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a commentServer) GetCommentByPartyUser(ctx echo.Context) error {
	pId := ctx.Param("pId")
	uId := ctx.Param("uId")

	c, err := a.cs.GetByPartyUser(ctx.Request().Context(), pId, uId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, c)
}
