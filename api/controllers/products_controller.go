/*
 * File: products_controller.go
 * Project: controllers
 * File Created: Thursday, 28th October 2021 10:08:21 am
 * Author: Ananda Yudhistira (anandabayu12@gmail.com)
 * -----
 * Last Modified: Thursday, 28th October 2021 12:33:04 pm
 * Modified By: Ananda Yudhistira (anandabayu12@gmail.com>)
 * -----
 * Copyright 2021 Ananda Yudhistira, -
 */
package controllers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anandabayu/sagara/api/auth"
	"github.com/anandabayu/sagara/api/models"
	"github.com/anandabayu/sagara/api/responses"
	"github.com/anandabayu/sagara/api/utils/formaterror"
	"github.com/gorilla/mux"
)

func (server *Server) CreateProduct(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	Product := models.Product{}
	err = json.Unmarshal(body, &Product)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	Product.Prepare()

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	Product.AuthorID = uid

	err = Product.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	fn := time.Now().Unix()

	coI := strings.Index(string(Product.Image), ",")
	sp := strings.Split(Product.Image, ",")

	switch strings.TrimSuffix(Product.Image[5:coI], ";base64") {
	case "image/png":
		fileName, err := base64toImage(sp[1], fmt.Sprintf("%d", fn), "png")
		if err != nil {
			log.Fatal(err)
		}
		Product.Image = fileName
	case "image/jpeg":
		fileName, err := base64toImage(sp[1], fmt.Sprintf("%d", fn), "jpg")
		if err != nil {
			log.Fatal(err)
		}
		Product.Image = fileName
	}

	ProductCreated, err := Product.SaveProduct(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	ProductCreated.Image = fmt.Sprintf("http://%s/%s", r.Host, ProductCreated.Image)
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, ProductCreated.ID))
	responses.JSON(w, http.StatusCreated, ProductCreated)
}

func (server *Server) GetProducts(w http.ResponseWriter, r *http.Request) {

	Product := models.Product{}

	Products, err := Product.FindAllProducts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	for i, p := range *Products {
		(*Products)[i].Image = fmt.Sprintf("http://%s/%s", r.Host, p.Image)
	}

	responses.JSON(w, http.StatusOK, Products)
}

func (server *Server) GetProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	Product := models.Product{}

	ProductReceived, err := Product.FindProductByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	ProductReceived.Image = fmt.Sprintf("http://%s/%s", r.Host, ProductReceived.Image)
	responses.JSON(w, http.StatusOK, ProductReceived)
}

func (server *Server) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the Product id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	fmt.Println(err)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the Product exist
	Product := models.Product{}
	err = server.DB.Debug().Model(models.Product{}).Where("id = ?", pid).Take(&Product).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Product not found"))
		return
	}

	// If a user attempt to update a Product not belonging to him
	fmt.Println(Product.AuthorID)
	fmt.Println(uid)
	if uid != Product.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	ProductUpdate := models.Product{}
	err = json.Unmarshal(body, &ProductUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	ProductUpdate.Prepare()
	ProductUpdate.AuthorID = uid
	err = ProductUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	fn := time.Now().Unix()

	coI := strings.Index(string(ProductUpdate.Image), ",")
	sp := strings.Split(ProductUpdate.Image, ",")

	switch strings.TrimSuffix(ProductUpdate.Image[5:coI], ";base64") {
	case "image/png":
		fileName, err := base64toImage(sp[1], fmt.Sprintf("%d", fn), "png")
		if err != nil {
			log.Fatal(err)
		}
		ProductUpdate.Image = fileName
	case "image/jpeg":
		fileName, err := base64toImage(sp[1], fmt.Sprintf("%d", fn), "jpg")
		if err != nil {
			log.Fatal(err)
		}
		ProductUpdate.Image = fileName
	}

	ProductUpdate.ID = Product.ID //this is important to tell the model the Product id to update, the other update field are set above

	ProductUpdated, err := ProductUpdate.UpdateAProduct(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	ProductUpdated.Image = fmt.Sprintf("http://%s/%s", r.Host, ProductUpdated.Image)
	responses.JSON(w, http.StatusOK, ProductUpdated)
}

func (server *Server) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid Product id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the Product exist
	Product := models.Product{}
	err = server.DB.Debug().Model(models.Product{}).Where("id = ?", pid).Take(&Product).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this Product?
	if uid != Product.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = Product.DeleteAProduct(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp := make(map[string]string)

	resp["message"] = "Product deleted successfully"

	responses.JSON(w, http.StatusOK, resp)
}

func base64toImage(data string, name string, types string) (string, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, formatString, err := image.Decode(reader)
	if err != nil {
		return "", err
	}
	bounds := m.Bounds()
	fmt.Println(bounds, formatString)

	//Encode from image format to writer
	fileName := fmt.Sprintf("images/%s.png", name)

	if types == "jpg" {
		fileName = fmt.Sprintf("images/%s.jpg", name)
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", err
	}

	switch types {
	case "png":
		err = png.Encode(f, m)
		if err != nil {
			// log.Fatal(err)
			return "", err
		}
	case "jpg":
		err = jpeg.Encode(f, m, &jpeg.Options{Quality: 75})
		if err != nil {
			return "", err
		}
	}

	fmt.Println("File", fileName, "created")
	return fileName, nil
}

func (server *Server) GetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	image := fmt.Sprintf("images/%s", vars["image"])

	fileBytes, err := ioutil.ReadFile(image)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	return
}
