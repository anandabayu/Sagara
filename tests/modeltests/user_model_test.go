/*
 * File: user_model_test.go
 * Project: modeltests
 * File Created: Thursday, 28th October 2021 12:42:59 pm
 * Author: Ananda Yudhistira (anandabayu12@gmail.com)
 * -----
 * Last Modified: Thursday, 28th October 2021 12:44:14 pm
 * Modified By: Ananda Yudhistira (anandabayu12@gmail.com>)
 * -----
 * Copyright 2021 Ananda Yudhistira, -
 */
package modeltests

import (
	"log"
	"testing"

	"github.com/anandabayu/sagara/api/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/go-playground/assert.v1"
)

func TestSaveUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	newUser := models.User{
		ID:       1,
		Email:    "test@gmail.com",
		Nickname: "test",
		Password: "password",
	}
	savedUser, err := newUser.SaveUser(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}
	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Email, savedUser.Email)
	assert.Equal(t, newUser.Nickname, savedUser.Nickname)
}
