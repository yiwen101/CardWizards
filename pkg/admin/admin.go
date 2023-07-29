package admin

import (
	"context"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/yiwen101/CardWizards/pkg/store"
)

/*
Offering a set of admin api for the runtime modification and control of the api gateway

By the time of writing, the admin package offer 16 apis. The details of 14 of them can be found in
in the rest of the files under the admin folder.

This file only adds two primitive api and one primitive authentication mechanism. Namely:
1 check whether proxy is on: GET /admin/proxy
2 turn on/off proxy: PUT /admin/proxy          json body: bool

The authentification is off by default. Can run the built app with -pdw flag to assign password and turn it on.
After turning it on, the "/admin" will only serve request with the correct "Password: ${password}" header.

Can use "AddRegister" function in other files to add more apis to the admin service.
*/

var registerlist []AdminRegister = []AdminRegister{registerProxy}

func Register(r *server.Hertz) {
	admin := r.Group("/admin")

	for _, f := range registerlist {
		f(admin)
	}
}

type AdminRegister func(*route.RouterGroup)

func AddRegister(f AdminRegister) {
	registerlist = append(registerlist, f)
}

func registerProxy(admin *route.RouterGroup) {
	password := store.InfoStore.Password
	if password != "" {
		admin.Use(func(ctx context.Context, c *app.RequestContext) {
			if string(c.GetHeader("Password")) != password {
				c.AbortWithMsg("wrong password", http.StatusBadRequest)
			}
		})
	}

	admin.GET("/proxy",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.CheckProxyStatus()
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, s)
		})
	admin.PUT("/proxy",
		func(ctx context.Context, c *app.RequestContext) {
			bytes, err := c.Body()
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			var b bool
			err = sonic.Unmarshal(bytes, &b)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			if b {
				err = store.InfoStore.TurnOnProxy()
			} else {
				err = store.InfoStore.TurnOffProxy()
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, "Proxy status updated")
		})
}
