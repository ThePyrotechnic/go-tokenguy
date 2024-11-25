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
	"encoding/json"
	"github.com/gin-gonic/gin/binding"
	"maps"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

var (
	router      *gin.Engine = nil
	initialized             = false
	publicKeys  map[string]*rsa.PublicKey
	privateKeys map[string]*rsa.PrivateKey
	digests     []byte
)

type TokenWrapper struct {
	Token string `json:"token" binding:"required"`
}

type TokenSigningWrapper struct {
	Claims string `json:"claims" binding:"required"`
	KID    string `json:"kid" binding:"required"`
}

func Router(_publicKeys map[string]*rsa.PublicKey, _privateKeys map[string]*rsa.PrivateKey) *gin.Engine {
	if !initialized {
		router = gin.Default()
		router.SetTrustedProxies(nil)
		router.POST("/validate", postValidate)
		router.POST("/sign", postSign)
		router.GET("/keys", getPublicKIDs)
		publicKeys = _publicKeys
		privateKeys = _privateKeys

		var err error
		digests, err = json.Marshal(gin.H{"keys": slices.Collect(maps.Keys(privateKeys))})
		if err != nil {
			panic(err)
		}
		initialized = true
	}

	return router
}

func getPublicKIDs(c *gin.Context) {
	c.Data(200, "application/json", digests)
}

func postSign(c *gin.Context) {
	var req TokenSigningWrapper
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	key, ok := privateKeys[req.KID]
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}
	signedToken, err := Sign(req.KID, key, req.Claims)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": signedToken})
}

func postValidate(c *gin.Context) {
	var token TokenWrapper
	if err := c.ShouldBindWith(&token, binding.JSON); err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	if Validate(publicKeys, token.Token) {
		c.JSON(http.StatusOK, gin.H{"valid": "true"})
	} else {
		c.Status(http.StatusForbidden)
	}
}
