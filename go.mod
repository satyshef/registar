module github.com/satyshef/registar

go 1.18

//replace github.com/satyshef/tdlib => ../tdlib

//replace github.com/satyshef/tdbot => ../tdbot

//replace github.com/satyshef/tdbot/chat => ../../telegram/tdbot/chat

require (
	github.com/BurntSushi/toml v1.1.0
	github.com/satyshef/go-tdlib v0.3.12
	github.com/satyshef/tdbot v0.3.0
	github.com/valyala/fasthttp v1.37.0
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/klauspost/compress v1.15.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
)
