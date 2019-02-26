// Copyright (C) 2018 O.S. Systems Sofware LTDA
//
// SPDX-License-Identifier: MIT

package webapi

import (
	"net/http"

	"github.com/asdine/storm"
	"github.com/gustavosbarreto/uhmicro/server/models"
	"github.com/labstack/echo"
)

const (
	GetAllProductsUrl = "/products"
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

	return c.JSON(http.StatusOK, products)
}
