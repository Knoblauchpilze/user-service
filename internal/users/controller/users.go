package controller

import (
	"net/http"

	"github.com/KnoblauchPilze/user-service/internal/users/service"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UserEndpoints(service service.UserService) rest.Routes {
	var out rest.Routes

	postHandler := fromUserServiceAwareHttpHandler(createUser, service)
	post := rest.NewRoute(http.MethodPost, false, "/users", postHandler)
	out = append(out, post)

	getHandler := fromUserServiceAwareHttpHandler(getUser, service)
	get := rest.NewResourceRoute(http.MethodGet, true, "/users", getHandler)
	out = append(out, get)

	listHandler := fromUserServiceAwareHttpHandler(listUsers, service)
	list := rest.NewRoute(http.MethodGet, true, "/users", listHandler)
	out = append(out, list)

	updateHandler := fromUserServiceAwareHttpHandler(updateUser, service)
	update := rest.NewResourceRoute(http.MethodPatch, true, "/users", updateHandler)
	out = append(out, update)

	deleteHandler := fromUserServiceAwareHttpHandler(deleteUser, service)
	delete := rest.NewResourceRoute(http.MethodDelete, true, "/users", deleteHandler)
	out = append(out, delete)

	loginByEmailHandler := fromUserServiceAwareHttpHandler(loginUserByEmail, service)
	loginByEmail := rest.NewRoute(http.MethodPost, false, "/users/sessions", loginByEmailHandler)
	out = append(out, loginByEmail)

	loginByIdHandler := fromUserServiceAwareHttpHandler(loginUserById, service)
	loginById := rest.NewResourceRoute(http.MethodPost, false, "/users/sessions", loginByIdHandler)
	out = append(out, loginById)

	logoutHandler := fromUserServiceAwareHttpHandler(logoutUser, service)
	logout := rest.NewResourceRoute(http.MethodDelete, true, "/users/sessions", logoutHandler)
	out = append(out, logout)

	return out
}

func createUser(c echo.Context, us service.UserService) error {
	// https://echo.labstack.com/docs/binding
	var userDtoRequest communication.UserDtoRequest
	err := c.Bind(&userDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user syntax")
	}

	out, err := us.Create(c.Request().Context(), userDtoRequest)
	if err != nil {
		if errors.IsErrorWithCode(err, db.DuplicatedKeySqlKey) {
			return c.JSON(http.StatusConflict, "Email already used")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, out)
}

func getUser(c echo.Context, us service.UserService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	out, err := us.Get(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func listUsers(c echo.Context, us service.UserService) error {
	out, err := us.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func updateUser(c echo.Context, us service.UserService) error {
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

	out, err := us.Update(c.Request().Context(), id, userDtoRequest)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		if errors.IsErrorWithCode(err, db.OptimisticLockException) {
			return c.JSON(http.StatusConflict, "User is not up to date")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func deleteUser(c echo.Context, us service.UserService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	err = us.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func loginUserByEmail(c echo.Context, us service.UserService) error {
	var userDtoRequest communication.UserDtoRequest
	err := c.Bind(&userDtoRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user syntax")
	}

	out, err := us.Login(c.Request().Context(), userDtoRequest)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}
		if errors.IsErrorWithCode(err, service.InvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, "Invalid credentials")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, out)
}

func loginUserById(c echo.Context, us service.UserService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	out, err := us.LoginById(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, out)
}

func logoutUser(c echo.Context, us service.UserService) error {
	maybeId := c.Param("id")
	id, err := uuid.Parse(maybeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id syntax")
	}

	err = us.Logout(c.Request().Context(), id)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingSqlRows) {
			return c.JSON(http.StatusNotFound, "No such user")
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
