# Token Guy

[![Go](https://github.com/ThePyrotechnic/go-tokenguy/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/ThePyrotechnic/go-tokenguy/actions/workflows/go.yml)

[![Docker Publish](https://github.com/ThePyrotechnic/go-tokenguy/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/ThePyrotechnic/go-tokenguy/actions/workflows/docker-publish.yml)


## Building

`make build`

The binary will be in `/bin/`


## Running

Place all desired public RS256 keys in a directory

Optionally, place matching RS256 private keys in a separate directory

`tokenguy server start --public-keys <public dir> [--private-keys <private dir>]`

See also `tokenguy server start -h`


### TLS

For production usage, serve the binary via a reverse proxy (like nginx) and set
`GIN_MODE=release`


## Usage as a web server

### Validating existing tokens

Send a POST request to `http://<host>:<port>/validate`

The body should look like: `{"token": "<token>"}`

Set `Content-Type: application/json` in your request

The response will be one of the following:

1. (status code `5XX`) if something is wrong with the server
2. (status code `403`) if the token or your request is invalid for any reason 
3. (status code `200`) `{"valid": true}` if the token is valid

### Signing new tokens

Send a GET request to `http://<host>:<port>/keys`

The response will be `{"keys": [<kid>, <kid>, ...]}`, where the `[<kid>]` array is a list of hashes of the public key
for each private key in the `--private-keys` directory. Pick one `<kid>`. 

Send a POST request to `http://<host>:<port>/sign`

Set `Content-Type: application/json` in your request


The body should look like: `{"kid": "<kid>", "claims": {"claim1": "value", ...}}`

The response will be `{"token": "<token>"}` where `<token>` is a valid JWT with your provided claims, signed
by the private key corresponding to the public `<kid>` hash you requested

Note that JWTs do not expire by default. If you would like your JWT to expire, add an `exp` to your claims.
(See also [the full list of "registered claims"](https://datatracker.ietf.org/doc/html/rfc7519#section-4.1))


## Usage from the CLI

`tokenguy validate --public-keys <dir> [token]` will check that `token` has been signed by one of the keys in `<dir>`


## License and Copyright
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
