package handler

import (
	"goyavision/internal/api/dto"
	"goyavision/internal/app"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterAlgorithm(g *echo.Group, d Deps) {
	svc := app.NewAlgorithmService(d.Repo)
	h := algorithmHandler{svc: svc}
	g.GET("/algorithms", h.List)
	g.POST("/algorithms", h.Create)
	g.GET("/algorithms/:id", h.Get)
	g.PUT("/algorithms/:id", h.Update)
	g.DELETE("/algorithms/:id", h.Delete)
}

type algorithmHandler struct {
	svc *app.AlgorithmService
}

func (h *algorithmHandler) List(c echo.Context) error {
	algorithms, err := h.svc.List(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(200, dto.AlgorithmsToResponse(algorithms))
}

func (h *algorithmHandler) Create(c echo.Context) error {
	var req dto.AlgorithmCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	alg, err := req.ToDomain()
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
	}

	result, err := h.svc.Create(c.Request().Context(), alg)
	if err != nil {
		return err
	}

	return c.JSON(201, dto.AlgorithmToResponse(result))
}

func (h *algorithmHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid algorithm id",
		})
	}

	alg, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.AlgorithmToResponse(alg))
}

func (h *algorithmHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid algorithm id",
		})
	}

	var req dto.AlgorithmUpdateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	alg := &dto.AlgorithmCreateReq{}
	if req.Name != nil {
		alg.Name = *req.Name
	}
	if req.Endpoint != nil {
		alg.Endpoint = *req.Endpoint
	}
	if req.InputSpec != nil {
		alg.InputSpec = req.InputSpec
	}
	if req.OutputSpec != nil {
		alg.OutputSpec = req.OutputSpec
	}

	domainAlg, err := alg.ToDomain()
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
	}

	result, err := h.svc.Update(c.Request().Context(), id, domainAlg)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.AlgorithmToResponse(result))
}

func (h *algorithmHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid algorithm id",
		})
	}

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(204)
}
