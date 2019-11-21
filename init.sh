docker pull golang
docker pull mongo
docker build -t go_react ./
mkdir -p mongodata/
mkdir -p mongodata/db
mkdir -p mongodata/configdb
