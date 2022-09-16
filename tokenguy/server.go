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
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
)

var (
	router      = gin.Default()
	initialized = false
	keys        map[string]*rsa.PublicKey
)

type TokenWrapper struct {
	Token string
}

func Router(_keys map[string]*rsa.PublicKey) *gin.Engine {
	if !initialized {
		router.SetTrustedProxies(nil)
		router.POST("/validate", postValidate)
		keys = _keys
		initialized = true
	}

	return router
}

func validate(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("Unable to parse KID")
		} else if key := keys[kid]; key != nil {
			return key, nil
		} else {
			return nil, fmt.Errorf("No matching KID")
		}
	}, jwt.WithValidMethods([]string{"RS256"}))

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Println(claims["name"], claims["admin"])
	} else {
		log.Println(err)
	}

	return token.Valid
}

func postValidate(c *gin.Context) {
	var token TokenWrapper
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	isValid := validate(token.Token)
	c.JSON(http.StatusOK, gin.H{"valid": isValid})
}
