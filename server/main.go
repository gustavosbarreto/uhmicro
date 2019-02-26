// Copyright (C) 2018 O.S. Systems Sofware LTDA
//
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/asdine/storm"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gustavosbarreto/uhmicro/server/api/agentapi"
	"github.com/gustavosbarreto/uhmicro/server/api/webapi"
	"github.com/gustavosbarreto/uhmicro/server/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cobra.OnInitialize(func() {
		viper.AutomaticEnv()
	})

	rootCmd := &cobra.Command{
		Use: "updatehub-ce-server",
		Run: execute,
	}

	rootCmd.PersistentFlags().StringP("db", "", "updatehub.db", "Database file")
	rootCmd.PersistentFlags().StringP("username", "", "admin", "Admin username")
	rootCmd.PersistentFlags().StringP("password", "", "admin", "Admin password")
	rootCmd.PersistentFlags().IntP("http", "", 8080, "HTTP listen address")
	rootCmd.PersistentFlags().StringP("dir", "", "./", "Packages storage dir")
	rootCmd.PersistentFlags().IntP("coap", "", 5683, "Coap server listen port")

	viper.BindPFlag("db", rootCmd.PersistentFlags().Lookup("db"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("http", rootCmd.PersistentFlags().Lookup("http"))
	viper.BindPFlag("dir", rootCmd.PersistentFlags().Lookup("dir"))
	viper.BindPFlag("coap", rootCmd.PersistentFlags().Lookup("coap"))

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func execute(cmd *cobra.Command, args []string) {
	db, err := storm.Open(viper.GetString("db"))
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.CORS())

	e.POST("/login", func(c echo.Context) error {
		var login struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		c.Bind(&login)

		if login.Username == "" {
			return echo.ErrUnauthorized
		}

		if login.Username == viper.GetString("username") && login.Password == viper.GetString("password") {
			token := jwt.New(jwt.SigningMethodHS256)

			claims := token.Claims.(jwt.MapClaims)
			claims["name"] = "root"
			claims["admin"] = true
			claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

			t, err := token.SignedString([]byte("secret"))
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, map[string]string{
				"token": t,
			})
		}

		return echo.ErrUnauthorized
	})

	agentApi := agentapi.NewAgentAPI(db)
	e.POST(agentapi.GetRolloutForDeviceUrl, agentApi.GetRolloutForDevice)
	e.POST(agentapi.ReportDeviceStateUrl, agentApi.ReportDeviceState)
	e.GET(agentapi.GetObjectFromPackageUrl, agentApi.GetObjectFromPackage)

	api := e.Group("/api")
	api.Use(middleware.JWT([]byte("secret")))

	devicesEndpoint := webapi.NewDevicesAPI(db)
	api.GET(webapi.GetAllDevicesUrl, devicesEndpoint.GetAllDevices)
	api.GET(webapi.GetDeviceUrl, devicesEndpoint.GetDevice)
	api.GET(webapi.GetDeviceRolloutReportsUrl, devicesEndpoint.GetDeviceRolloutReports)

	packagesEndpoint := webapi.NewPackagesAPI(db, viper.GetString("dir"))
	api.GET(webapi.GetAllPackagesUrl, packagesEndpoint.GetAllPackages)
	api.GET(webapi.GetPackageUrl, packagesEndpoint.GetPackage)
	api.POST(webapi.UploadPackageUrl, packagesEndpoint.UploadPackage)

	rolloutsEndpoint := webapi.NewRolloutsAPI(db)
	api.GET(webapi.GetAllRolloutsUrl, rolloutsEndpoint.GetAllRollouts)
	api.GET(webapi.GetRolloutUrl, rolloutsEndpoint.GetRollout)
	api.GET(webapi.GetRolloutStatisticsUrl, rolloutsEndpoint.GetRolloutStatistics)
	api.GET(webapi.GetRolloutDevicesUrl, rolloutsEndpoint.GetRolloutDevices)
	api.POST(webapi.CreateRolloutUrl, rolloutsEndpoint.CreateRollout)
	api.PUT(webapi.StopRolloutUrl, rolloutsEndpoint.StopRollout)

	namespacesEndpoint := webapi.NewNamespacesAPI(db)
	e.GET(webapi.GetAllNamespacesUrl, namespacesEndpoint.GetAllNamespaces)

	productsEndpoint := webapi.NewProductsAPI(db)
	e.GET(webapi.GetAllProductsUrl, productsEndpoint.GetAllProducts)

	n := models.Namespace{Name: "Roberto", UID: "namespace-uid-1"}
	err = db.Save(&n)
	fmt.Println(err)

	e.GET("/hydra/oauth2/auth", func(c echo.Context) error {
		redirect := c.QueryParam("redirect_uri")
		state := c.QueryParam("state")

		params := url.Values{
			"access_token": []string{"1234"},
			"state":        []string{state},
			"expires_in":   []string{"3600"},
		}

		dst, err := url.Parse(redirect)
		if err != nil {
			return err
		}

		dst.RawQuery = params.Encode()

		fmt.Println(dst.String())

		return c.Redirect(http.StatusMovedPermanently, dst.String())
	})

	e.GET("/me", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"name":      "Roberto Souza",
			"email":     "gustavosbarreto@gmail.com",
			"email_md5": "95a06375b611671eb47965e68de459d6",
		})
	})

	go func() {
		log.Fatal(e.Start(fmt.Sprintf(":%d", viper.GetInt("http"))))
	}()

	go func() {
		log.Fatal(startCoapServer(viper.GetInt("coap"), viper.GetInt("http")))
	}()

	select {}
}
