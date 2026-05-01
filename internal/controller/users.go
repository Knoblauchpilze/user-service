package controller

import (
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/pgx"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/user-service/internal/service"
	"github.com/Knoblauchpilze/user-service/pkg/communication"
	"github.com/Knoblauchpilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func UserEndpoints(service service.UserService) rest.Routes {
	var out rest.Routes

	postHandler := createServiceAwareHttpHandler(createUser, service)
	post := rest.NewRoute(http.MethodPost, "", postHandler)
	out = append(out, post)

	getHandler := createServiceAwareHttpHandler(getUser, service)
	get := rest.NewRoute(http.MethodGet, ":id", getHandler)
	out = append(out, get)

	listHandler := createServiceAwareHttpHandler(listUsers, service)
	list := rest.NewRoute(http.MethodGet, "", listHandler)
	out = append(out, list)

	updateHandler := createServiceAwareHttpHandler(updateUser, service)
	update := rest.NewRoute(http.MethodPatch, ":id", updateHandler)
	out = append(out, update)

	deleteHandler := createServiceAwareHttpHandler(deleteUser, service)
	delete := rest.NewRoute(http.MethodDelete, ":id", deleteHandler)
	out = append(out, delete)

	loginByEmailHandler := createServiceAwareHttpHandler(loginUserByEmail, service)
	loginByEmail := rest.NewRoute(http.MethodPost, "/sessions", loginByEmailHandler)
	out = append(out, loginByEmail)

	logoutHandler := createServiceAwareHttpHandler(logoutUser, service)
	logout := rest.NewRoute(http.MethodDelete, "/sessions/:id", logoutHandler)
	out = append(out, logout)

	return out
}

// createUser godoc
//
// @Summary Create user
// @Description Creates a user from the provided credentials.
// @Tags users
// @Produce json
// @Param user body communication.UserDtoRequest true "User payload"
// @Success 201 {object} rest.ResponseEnvelope[communication.UserDtoResponse]
// @Failure 400 {object} rest.ResponseEnvelope[string] "Invalid user syntax, email, or password"
// @Failure 409 {object} rest.ResponseEnvelope[string] "Email already in use"
// @Failure 500 {object} rest.ResponseEnvelope[string] "Internal server error"
// @Router /users [post]
func createUser(c *echo.Context, s service.UserService) error {
	// https://echo.labstack.com/docs/binding
	var userDtoRequest communication.UserDtoRequest
	err := c.Bind(&userDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user syntax")
	}

	out, err := s.Create(c.Request().Context(), userDtoRequest)
	if err != nil {
		if errors.IsErrorWithCode(err, service.InvalidEmail) {
			return c.JSON(http.StatusBadRequest, "Invalid email")
		}
		if errors.IsErrorWithCode(err, service.InvalidPassword) {
			return c.JSON(http.StatusBadRequest, "Invalid password")
		}
		if errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation) {
			return c.JSON(http.StatusConflict, "Email already in use")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, out)
}

// getUser godoc
//
// @Summary Get user
// @Description Returns a user by its identifier.
// @Tags users
// @Produce json
// @Param id path string true "User ID" Format(uuid)
// @Success 200 {object} rest.ResponseEnvelope[communication.UserDtoResponse]
// @Failure 400 {object} rest.ResponseEnvelope[string] "Invalid id syntax"
// @Failure 404 {object} rest.ResponseEnvelope[string] "No such user"
// @Failure 500 {object} rest.ResponseEnvelope[string] "Internal server error"
// @Router /users/{id} [get]
func getUser(c *echo.Context, s service.UserService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	out, err := s.Get(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

// listUsers godoc
//
// @Summary List users
// @Description Returns the identifiers of all users.
// @Tags users
// @Produce json
// @Success 200 {object} rest.ResponseEnvelope[[]string]
// @Failure 500 {object} rest.ResponseEnvelope[string] "Internal server error"
// @Router /users [get]
func listUsers(c *echo.Context, s service.UserService) error {
	out, err := s.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

// updateUser godoc
//
// @Summary Update user
// @Description Updates a user identified by its identifier.
// @Tags users
// @Produce json
// @Param id path string true "User ID" Format(uuid)
// @Param user body communication.UserDtoRequest true "User payload"
// @Success 200 {object} rest.ResponseEnvelope[communication.UserDtoResponse]
// @Failure 400 {object} rest.ResponseEnvelope[string] "Invalid id or user syntax"
// @Failure 404 {object} rest.ResponseEnvelope[string] "No such user"
// @Failure 409 {object} rest.ResponseEnvelope[string] "User is not up to date"
// @Failure 500 {object} rest.ResponseEnvelope[string] "Internal server error"
// @Router /users/{id} [patch]
func updateUser(c *echo.Context, s service.UserService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	var userDtoRequest communication.UserDtoRequest
	err = c.Bind(&userDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user syntax")
	}

	out, err := s.Update(c.Request().Context(), id, userDtoRequest)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		if errors.IsErrorWithCode(err, repositories.OptimisticLockException) {
			return c.JSON(http.StatusConflict, "User is not up to date")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

// deleteUser godoc
//
// @Summary Delete user
// @Description Deletes a user identified by its identifier.
// @Tags users
// @Param id path string true "User ID" Format(uuid)
// @Success 204
// @Failure 400 {object} rest.ResponseEnvelope[string] "Invalid id syntax"
// @Failure 500 {object} rest.ResponseEnvelope[string] "Internal server error"
// @Router /users/{id} [delete]
func deleteUser(c *echo.Context, s service.UserService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	err = s.Delete(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// loginUserByEmail godoc
//
// @Summary Create session
// @Description Authenticates a user with email and password and returns an API key.
// @Tags sessions
// @Produce json
// @Param user body communication.UserDtoRequest true "User credentials"
// @Success 201 {object} rest.ResponseEnvelope[communication.ApiKeyDtoResponse]
// @Failure 400 {object} rest.ResponseEnvelope[string] "Invalid user syntax"
// @Failure 401 {object} rest.ResponseEnvelope[string] "Invalid credentials"
// @Failure 404 {object} rest.ResponseEnvelope[string] "No such user"
// @Failure 500 {object} rest.ResponseEnvelope[string] "Internal server error"
// @Router /users/sessions [post]
func loginUserByEmail(c *echo.Context, s service.UserService) error {
	var userDtoRequest communication.UserDtoRequest
	err := c.Bind(&userDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user syntax")
	}

	out, err := s.Login(c.Request().Context(), userDtoRequest)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}
		if errors.IsErrorWithCode(err, service.InvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, "Invalid credentials")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, out)
}

// logoutUser godoc
//
// @Summary Delete session
// @Description Revokes the active session for the specified user.
// @Tags sessions
// @Param id path string true "User ID" Format(uuid)
// @Success 204
// @Failure 400 {object} rest.ResponseEnvelope[string] "Invalid id syntax"
// @Failure 404 {object} rest.ResponseEnvelope[string] "No such user"
// @Failure 500 {object} rest.ResponseEnvelope[string] "Internal server error"
// @Router /users/sessions/{id} [delete]
func logoutUser(c *echo.Context, s service.UserService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	err = s.Logout(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
