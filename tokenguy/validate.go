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
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

func GetKeys() map[string]*rsa.PublicKey {
	keysDir := viper.GetString("server.keys")
	fileinfo, err := os.Stat(keysDir)
	if err != nil {
		panic(err)
	}
	if !fileinfo.IsDir() {
		panic(fmt.Errorf("Provided keys directory is not a directory"))
	}
	matches, err := filepath.Glob(filepath.Join(keysDir, "*"))
	if err != nil {
		panic(err)
	}
	if matches == nil {
		panic(fmt.Errorf("Provided keys directory is empty"))
	}

	keyMap := make(map[string]*rsa.PublicKey)
	for a := 0; a < len(matches); a++ {
		data, err := os.ReadFile(matches[a])
		if err != nil {
			log.Println(matches[a], ": ", err)
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM(data)
		if err != nil {
			log.Println(matches[a], ": ", err)
		}
		keyMap[filepath.Base(matches[a])] = key
	}

	return keyMap
}

func Validate(keys map[string]*rsa.PublicKey, tokenString string) bool {
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
