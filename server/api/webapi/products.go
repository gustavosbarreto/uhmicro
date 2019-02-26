// Copyright (C) 2018 O.S. Systems Sofware LTDA
//
// SPDX-License-Identifier: MIT

package webapi

import (
	"net/http"

	"github.com/asdine/storm"
	"github.com/google/uuid"
	"github.com/gustavosbarreto/uhmicro/server/models"
	"github.com/labstack/echo"
)

const (
	GetAllProductsUrl = "/products"
	CreateProductUrl  = "/products"
)

type ProductsAPI struct {
	db *storm.DB
}

func NewProductsAPI(db *storm.DB) *ProductsAPI {
	return &ProductsAPI{db: db}
}

func (api *ProductsAPI) GetAllProducts(c echo.Context) error {
	var products []models.Product
	if err := api.db.All(&products); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"products": products,
	})
}

func (api *ProductsAPI) CreateProduct(c echo.Context) error {
	product := models.Product{}

	if err := c.Bind(&product); err != nil {
		return err
	}

	product.UID = uuid.New().String()

	if err := api.db.Save(&product); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, product)
}
