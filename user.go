package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserEntity struct {
	Id       int    `json:"id,omitempty"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Fullname string `json:"fullname,omitempty"`
}

type User struct {
	UserEntity
	Links `json:"links,omitempty"`
}

//------------------------------------------------
// REST-API
//------------------------------------------------

func findAllUsers(c *gin.Context) {
	users, err := getAllUsers()
	if err != nil {
		renderError(c, http.StatusInternalServerError, err.Error())
		return
	}
	ua := make([]*User, 0)
	for _, ue := range users {
		lb := New(c.Request)
		lb.Add("self", fmt.Sprintf("/api/user/%d", ue.Id))
		u := User{
			*ue,
			lb.Links,
		}
		ua = append(ua, &u)
	}
	c.JSON(http.StatusOK, ua)
}

func findUserById(c *gin.Context) {
	s := c.Param("id")
	id, err := strconv.Atoi(s)
	if err != nil {
		renderError(c, http.StatusBadRequest, "id must be an int")
		return
	}
	ue, err := getUserById(id)
	if err != nil {
		renderError(c, http.StatusNotFound, err.Error())
		return
	}
	lb := New(c.Request)
	lb.Add("self", fmt.Sprintf("/api/user/%d", ue.Id))
	user := User{
		*ue,
		lb.Links,
	}
	c.JSON(http.StatusOK, user)
}

func addUser(c *gin.Context) {
	u := UserEntity{}
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
	lb := New(c.Request)
	lb.Add("self", fmt.Sprintf("/api/user/%d", u2.Id))
	user := User{
		*u2,
		lb.Links,
	}
	c.JSON(http.StatusOK, user)
}

func updateUser(c *gin.Context) {
	u := UserEntity{}
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
	lb := New(c.Request)
	lb.Add("self", fmt.Sprintf("/api/user/%d", u2.Id))
	user := User{
		*u2,
		lb.Links,
	}
	c.JSON(http.StatusOK, user)
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

var users = make([]*UserEntity, 0)

// Add user to user list
func create(u *UserEntity) (*UserEntity, error) {
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

func delete(u *UserEntity) error {
	for i, us := range users {
		if us.Id == u.Id {
			users = append(users[:i], users[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("user %s with id %d not found", u.UserName, u.Id)
}

func update(u *UserEntity) (*UserEntity, error) {
	us, err := getUserById(u.Id)
	if err != nil {
		return nil, err
	}
	us.UserName = u.UserName
	us.Email = u.Email
	us.Fullname = u.Fullname
	return u, nil
}

func getUserById(id int) (*UserEntity, error) {
	for _, us := range users {
		if us.Id == id {
			return us, nil
		}
	}
	return nil, fmt.Errorf("user with id %d not found", id)
}

func findUserByUsername(username string) (*UserEntity, error) {
	for _, us := range users {
		if us.UserName == username {
			return us, nil
		}
	}
	return nil, fmt.Errorf("user with username %s not found", username)
}

func getAllUsers() ([]*UserEntity, error) {
	return users, nil
}
