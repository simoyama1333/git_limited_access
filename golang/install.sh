go get -u github.com/labstack/echo/...
go get -u github.com/golang/dep/cmd/dep
dep init
dep ensure -add "go.mongodb.org/mongo-driver/mongo@~1.0.0"