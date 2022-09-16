# Token Guy

## Building

`make build`

The binary will be in `/bin/`


## Running

Place all desired validation keys in a directory

`tokenguy server start --keys-directory <dir>`

See `tokenguy server start -h` for more options

### TLS

For production usage, serve the binary via a reverse proxy (like nginx)


## Usage

Send a POST request to `http://<host>:<port>/validate`

The body should look like: `{"token": "<token>"}`

Set `Content-Type: application/json` in your request

The response will be one of the following:

1. (status code `4XX`) `{"error": "<message>"}` if something is wrong with your request
2. (status code `5XX`) if something is wrong with the server
3. (status code `200`) `{"valid": true/false}` if the token is valid or invalid


## License and Copyright
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
