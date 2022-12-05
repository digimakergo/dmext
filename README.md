# dmext
Extended utilities for digimaker, eg. outlook email, log, callbacks, etc

## Outlook email
under `dm.yaml/general`
```yaml
send_from: ""
mail_host: ""
mail_password: ""
```

Import mail in `main.go`
```go
	_ "github.com/digimakergo/dmext/mail"
```

## DMEditor
Configure dmeditor server `dm.yaml/general`. For running dmeditor server please check https://github.com/digimakergo/dmeditor
```yaml
  dmeditor_server_url: "http://localhost:8086/dmeditor"
```


Import dmeditor in `main.go`
```go
	_ "github.com/digimakergo/dmext/dmeditor"
```