# Idea

Idea behind go mod is to send files between computers on local lan network. First thing is to create cli app which will from one side

- Serve file from my computer
- Download file from someone computer by ip
- Create QR code which will allow to download file from website

## CLI Commands

### serve

`serve` command will allow to share file with other person with simple.

> godrop serve {{file}}||{{directory}}

- **--qr=true** - will allow to create qr code to share it, if qr code enabled no password is possible
- **--once** - will close application after first succesful download
- **--password** - password in header (hashed before sending)
- **--n** - n - times allow to download
- **--port** - port
  !! For directory we will zip it first
  We will then display our current ip address and waiting state

### get

`get` command will allow to get this file on other pc

> godrop get --from {{IP_ADDRESS}}

- **--from** - ip address
- **--password** - password

this will allow to download file

## Architecture

Internal logic is divided into server part, client and qr code generation

### Server

```go
// public api
type ServerOptions struct{
	password *string // default nil
	qrEnabled bool // default false
	nTimes *int// default MAX_INT
	port *int
}
// internal
type serverConfig struct {
	password *string
	qrEnabled bool
	nTimes int
	port int
}
```

### Client

```go

type ClientOptions struct{
	from *string // required
	password *string // default nil
}

### QR
```
