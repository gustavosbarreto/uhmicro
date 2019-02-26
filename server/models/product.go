// Copyright (C) 2018 O.S. Systems Sofware LTDA
//
// SPDX-License-Identifier: MIT

package models

type Product struct {
	UID          string `storm:"id" json:"uid"`
	Name         string `json:"name"`
	Organization string `storm:"index" json:"organization"`
}
