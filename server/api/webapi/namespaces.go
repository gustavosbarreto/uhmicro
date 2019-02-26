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
	GetAllNamespacesUrl = "/namespaces"
)

type NamespacesAPI struct {
	db *storm.DB
}

func NewNamespacesAPI(db *storm.DB) *NamespacesAPI {
	return &NamespacesAPI{db: db}
}

func (api *NamespacesAPI) GetAllNamespaces(c echo.Context) error {
	var namespaces []models.Namespace
	if err := api.db.All(&namespaces); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, namespaces)
}
