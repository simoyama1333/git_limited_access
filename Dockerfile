FROM golang
RUN apt-get update
RUN apt-get install -y nginx
ADD default.conf /etc/nginx/conf.d/
