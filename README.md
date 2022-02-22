# coursebench-backend

course-bench 的后端服务

## Build

* Install Go
* Install Docker
* Install Docker-Compose

```shell
git clone git@github.com:ShanghaitechGeekPie/coursebench-backend.git
cd coursebench-backend
go build -o bin/coursebench-backend cmd/coursebench-backend/main.go
go build -o bin/import_course cmd/import_course/main.go
```

## Configure

```shell
cp config.json.example config.json
```

Edit ``config.json`` as you like.

Edit files in `build` directory as you like.

## Import data (optional)

```shell
./bin/import_course
```

## Run

```shell
cd build
sudo docker-compose up -d
cd ..
./bin/coursebench-backend
```