module crawl_xsky_job_info

go 1.17

require github.com/beego/beego/v2 v2.0.2

require (
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
)

replace (
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v4 v4.3.0
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.5.0
	github.com/miekg/dns => github.com/miekg/dns v1.1.46
)
