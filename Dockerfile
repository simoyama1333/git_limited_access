FROM golang
RUN apt-get update
RUN apt-get install -y nginx
RUN go get -u github.com/labstack/echo/...
RUN go get -u github.com/golang/dep/cmd/dep