package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type User struct {
	Id       int    `json:"id,omitempty"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Fullname string `json:"fullname,omitempty"`
}

//------------------------------------------------
// REST-API
//------------------------------------------------

func findAllUsers(c *gin.Context) {
	u, err := getAllUsers()
	if err != nil {
		renderError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, u)
}

func findUserById(c *gin.Context) {
	s := c.Param("id")
	id, err := strconv.Atoi(s)
	if err != nil {
		renderError(c, http.StatusBadRequest, "id must be an int")
		return
	}
	u, err := getUserById(id)
	if err != nil {
		renderError(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, u)
}

func addUser(c *gin.Context) {
	u := User{}
	err := c.BindJSON(&u)
	if err != nil {
		renderError(c, http.StatusBadRequest, "error parsing user data. "+err.Error())
		return
	}
	if u.Email == "" || u.UserName == "" {
		renderError(c, http.StatusBadRequest, "username and email must be provided")
		return
	}

	_, err = findUserByUsername(u.UserName)
	if err == nil {
		renderError(c, http.StatusBadRequest, "username must not exists")
		return
	}
	u2, err := create(&u)
	if err != nil {
		renderError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, u2)
}

func updateUser(c *gin.Context) {
	u := User{}
	err := c.BindJSON(&u)
	if err != nil {
		renderError(c, http.StatusBadRequest, "error parsing user data. "+err.Error())
		return
	}
	if u.Id == 0 || u.Email == "" || u.UserName == "" {
		renderError(c, http.StatusBadRequest, "id, username and email must be provided")
		return
	}

	u2, err := update(&u)
	if err != nil {
		renderError(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, u2)
}

func deleteUser(c *gin.Context) {
	s := c.Param("id")
	id, err := strconv.Atoi(s)
	if err != nil {
		renderError(c, http.StatusBadRequest, "id must be an int")
		return
	}
	u, err := getUserById(id)
	if err != nil {
		renderError(c, http.StatusNotFound, err.Error())
		return
	}
	err = delete(u)
	if err != nil {
		renderError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

//------------------------------------------------
// BACKEND
//------------------------------------------------

var users = make([]*User, 0)

// add user to user list
func create(u *User) (*User, error) {
	if u.Id != 0 {
		return nil, fmt.Errorf("id must not be set")
	}
	// id is highest id + 1
	for _, us := range users {
		if us.Id > u.Id {
			u.Id = us.Id
		}
	}
	u.Id = u.Id + 1
	users = append(users, u)
	return u, nil
}

func delete(u *User) error {
	for i, us := range users {
		if us.Id == u.Id {
			users = append(users[:i], users[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("user %s with id %d not found", u.UserName, u.Id)
}

func update(u *User) (*User, error) {
	us, err := getUserById(u.Id)
	if err != nil {
		return nil, err
	}
	us.UserName = u.UserName
	us.Email = u.Email
	us.Fullname = u.Fullname
	return u, nil
}

func getUserById(id int) (*User, error) {
	for _, us := range users {
		if us.Id == id {
			return us, nil
		}
	}
	return nil, fmt.Errorf("user with id %d not found", id)
}

func findUserByUsername(username string) (*User, error) {
	for _, us := range users {
		if us.UserName == username {
			return us, nil
		}
	}
	return nil, fmt.Errorf("user with username %s not found", username)
}

func getAllUsers() ([]*User, error) {
	return users, nil
}
