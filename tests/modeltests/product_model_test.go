/*
 * File: product_model_test.go
 * Project: modeltests
 * File Created: Thursday, 28th October 2021 12:50:32 pm
 * Author: Ananda Yudhistira (anandabayu12@gmail.com)
 * -----
 * Last Modified: Thursday, 28th October 2021 12:55:55 pm
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

func TestFindAllProducts(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}
	_, _, err = seedUsersAndProducts()
	if err != nil {
		log.Fatalf("Error seeding user and post  table %v\n", err)
	}
	posts, err := postInstance.FindAllProducts(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the posts: %v\n", err)
		return
	}
	assert.Equal(t, len(*posts), 2)
}

func TestSaveProduct(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatalf("Error user and post refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newProduct := models.Product{
		ID:       1,
		Title:    "This is the title",
		Content:  "This is the content",
		AuthorID: user.ID,
	}
	savedProduct, err := newProduct.SaveProduct(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the post: %v\n", err)
		return
	}
	assert.Equal(t, newProduct.ID, savedProduct.ID)
	assert.Equal(t, newProduct.Title, savedProduct.Title)
	assert.Equal(t, newProduct.Content, savedProduct.Content)
	assert.Equal(t, newProduct.AuthorID, savedProduct.AuthorID)

}

func TestGetProductByID(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	post, err := seedOneUserAndOneProduct()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundProduct, err := postInstance.FindProductByID(server.DB, post.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundProduct.ID, post.ID)
	assert.Equal(t, foundProduct.Title, post.Title)
	assert.Equal(t, foundProduct.Content, post.Content)
}

func TestUpdateAProduct(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	post, err := seedOneUserAndOneProduct()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	postUpdate := models.Product{
		ID:       1,
		Title:    "modiUpdate",
		Content:  "modiupdate@gmail.com",
		AuthorID: post.AuthorID,
	}
	updatedProduct, err := postUpdate.UpdateAProduct(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedProduct.ID, postUpdate.ID)
	assert.Equal(t, updatedProduct.Title, postUpdate.Title)
	assert.Equal(t, updatedProduct.Content, postUpdate.Content)
	assert.Equal(t, updatedProduct.AuthorID, postUpdate.AuthorID)
}

func TestDeleteAProduct(t *testing.T) {

	err := refreshUserAndProductTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	post, err := seedOneUserAndOneProduct()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := postInstance.DeleteAProduct(server.DB, post.ID, post.AuthorID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, int64(1))
}
