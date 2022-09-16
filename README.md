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
