/*
 * File: home_controller.go
 * Project: controllers
 * File Created: Thursday, 28th October 2021 10:06:36 am
 * Author: Ananda Yudhistira (anandabayu12@gmail.com)
 * -----
 * Last Modified: Thursday, 28th October 2021 10:07:07 am
 * Modified By: Ananda Yudhistira (anandabayu12@gmail.com>)
 * -----
 * Copyright 2021 Ananda Yudhistira, -
 */

package controllers

import (
	"net/http"

	"github.com/anandabayu/sagara/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Sagara Test")

}
