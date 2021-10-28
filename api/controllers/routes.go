/*
 * File: routes.go
 * Project: controllers
 * File Created: Thursday, 28th October 2021 10:17:31 am
 * Author: Ananda Yudhistira (anandabayu12@gmail.com)
 * -----
 * Last Modified: Thursday, 28th October 2021 12:08:10 pm
 * Modified By: Ananda Yudhistira (anandabayu12@gmail.com>)
 * -----
 * Copyright 2021 Ananda Yudhistira, -
 */

package controllers

import "github.com/anandabayu/sagara/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")

	//Products routes
	s.Router.HandleFunc("/products", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateProduct))).Methods("POST")
	s.Router.HandleFunc("/products", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetProducts))).Methods("GET")
	s.Router.HandleFunc("/products/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetProduct))).Methods("GET")
	s.Router.HandleFunc("/products/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateProduct))).Methods("PUT")
	s.Router.HandleFunc("/products/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.DeleteProduct))).Methods("DELETE")

	s.Router.HandleFunc("/images/{image}", s.GetImage).Methods("GET")
}
