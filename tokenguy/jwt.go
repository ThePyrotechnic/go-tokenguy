/*
tokenguy is a web server which validates JWTs
Copyright (C) 2024  Michael Manis

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
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func getKeys(dir string) ([]string, error) {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		panic(err)
	}
	if !fileInfo.IsDir() {
		return []string{}, fmt.Errorf("provided keys directory is not a directory")
	}
	matches, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return []string{}, err
	}
	if matches == nil {
		return []string{}, fmt.Errorf("provided keys directory is empty")
	}
	return matches, nil
}

func pubKeyDigest(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func GetPublicKeys() map[string]*rsa.PublicKey {
	keysDir := viper.GetString("server.keys.public")
	matches, err := getKeys(keysDir)
	if err != nil {
		panic(err)
	}

	keyMap := make(map[string]*rsa.PublicKey)
	for a := 0; a < len(matches); a++ {
		data, err := os.ReadFile(matches[a])
		if err != nil {
			log.Println(matches[a], ": ", err)
			continue
		}

		key, err := jwt.ParseRSAPublicKeyFromPEM(data)
		if err != nil {
			log.Println(matches[a], ": ", err)
			continue
		}
		hDigest := pubKeyDigest(key.N.Bytes())
		keyMap[hDigest] = key
	}

	return keyMap
}

func GetPrivateKeys() map[string]*rsa.PrivateKey {
	keysDir := viper.GetString("server.keys.private")
	matches, err := getKeys(keysDir)
	if err != nil {
		panic(err)
	}

	keyMap := make(map[string]*rsa.PrivateKey)
	for a := 0; a < len(matches); a++ {
		data, err := os.ReadFile(matches[a])
		if err != nil {
			log.Println(matches[a], ": ", err)
			continue
		}
		key, err := jwt.ParseRSAPrivateKeyFromPEM(data)
		if err != nil {
			log.Println(matches[a], ": ", err)
			continue
		}

		keyMap[pubKeyDigest(key.PublicKey.N.Bytes())] = key
	}

	return keyMap
}

func Sign(kid string, key *rsa.PrivateKey, claims string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Header["kid"] = kid

	h, err := json.Marshal(token.Header)
	if err != nil {
		return "", err
	}

	signingString := token.EncodeSegment(h) + "." + token.EncodeSegment([]byte(claims))

	sig, err := token.Method.Sign(signingString, key)
	if err != nil {
		return "", err
	}

	return signingString + "." + token.EncodeSegment(sig), nil
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

	if err != nil {
		return false
	}

	return token.Valid
}
