# Simple webhook for [AutoFAQ async requests](https://redocly.github.io/redoc/?url=https://app.swaggerhub.com/apiproxy/registry/AutoFAQ.ai/external-api/2.1.4#operation/set_webhook)

## Configuration
### by File
Use config.yaml file to configurate server. Where:
* `server`:
    * `address` - binding address. Default: `0.0.0.0`
    * `port` - binding port. Default: `8000`
* `client`
    * `url` - webhook client url. Server send all AutoFAQ messages to that client. Default: `http://127.0.0.1:3000`

Example:
```yaml
---
server:
  address: 192.168.1.2
  port: 8080
client:
    url: http://client_address/webhook
```

### by Environmant variables
`SERVER_ADDRESS` - binding address. Defaults to `server.address`
`SERVER_PORT` - binding port. Default: `server.port` 
`CLIENT_URL` - webhook client url. Server send all AutoFAQ messages to that client. Default: `client.url`

## ToDo
* Params from config file or env
* Prometheus metrics