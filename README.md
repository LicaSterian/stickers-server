# stickers-server

This webserver requires a local mongo instance running.

**ROUTES**

/stickers/list
/stickers/add
/sticker/:filename

You can test the routes with `curl`

**DEV**

`go run main.go`

**BUILD**

`go build`

**PRODUCTION**

Run the _stickers-server_ binary after build