package handler

import (
	"net/http"

	"goyavision/internal/api/dto"
	"goyavision/internal/app"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterUser(g *echo.Group, d Deps) {
	svc := app.NewUserService(d.Repo)
	h := userHandler{svc: svc}

	g.GET("/users", h.List)
	g.POST("/users", h.Create)
	g.GET("/users/:id", h.Get)
	g.PUT("/users/:id", h.Update)
	g.DELETE("/users/:id", h.Delete)
	g.POST("/users/:id/reset-password", h.ResetPassword)
}

type userHandler struct {
	svc *app.UserService
}

func (h *userHandler) List(c echo.Context) error {
	var query dto.UserListQuery
	if err := c.Bind(&query); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid query parameters",
		})
	}

	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	users, total, err := h.svc.List(c.Request().Context(), query.Status, query.Limit, query.Offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.UserListResponse{
		Items: dto.UsersToResponse(users),
		Total: total,
	})
}

func (h *userHandler) Create(c echo.Context) error {
	var req dto.UserCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "username and password are required",
		})
	}

	if len(req.Password) < 6 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "password must be at least 6 characters",
		})
	}

	user, err := h.svc.Create(c.Request().Context(), &app.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Status:   req.Status,
		RoleIDs:  req.RoleIDs,
	})
	if err != nil {
		if err == app.ErrUsernameExists {
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "Conflict",
				Message: "username already exists",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.UserToResponse(user))
}

func (h *userHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid user id",
		})
	}

	user, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		if err == app.ErrUserNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "user not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.UserToResponse(user))
}

func (h *userHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid user id",
		})
	}

	var req dto.UserUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	user, err := h.svc.Update(c.Request().Context(), id, &app.UpdateUserRequest{
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Status:   req.Status,
		Password: req.Password,
		RoleIDs:  req.RoleIDs,
	})
	if err != nil {
		if err == app.ErrUserNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "user not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.UserToResponse(user))
}

func (h *userHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid user id",
		})
	}

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		if err == app.ErrUserNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "user not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

type resetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

func (h *userHandler) ResetPassword(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid user id",
		})
	}

	var req resetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	if len(req.NewPassword) < 6 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "password must be at least 6 characters",
		})
	}

	if err := h.svc.ResetPassword(c.Request().Context(), id, req.NewPassword); err != nil {
		if err == app.ErrUserNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "user not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "password reset successfully",
	})
}
