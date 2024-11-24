/*
tokenguy is a web server which validates JWTs
Copyright (C) 2022  Michael Manis

	This program is free software; you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation; either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program; if not, write to the Free Software Foundation,
	Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301  USA
*/
package tokenguy

import (
	"crypto/rsa"
	"github.com/gin-gonic/gin/binding"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	router      *gin.Engine = nil
	initialized             = false
	keys        map[string]*rsa.PublicKey
)

type TokenWrapper struct {
	Token string `json:"token" binding:"required"`
}

func Router(_keys map[string]*rsa.PublicKey) *gin.Engine {
	if !initialized {
		router = gin.Default()
		router.SetTrustedProxies(nil)
		router.POST("/validate", postValidate)
		keys = _keys
		initialized = true
	}

	return router
}

func postValidate(c *gin.Context) {
	var token TokenWrapper
	if err := c.ShouldBindWith(&token, binding.JSON); err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	if Validate(keys, token.Token) {
		c.JSON(http.StatusOK, gin.H{"valid": "true"})
	} else {
		c.Status(http.StatusForbidden)
	}
}
