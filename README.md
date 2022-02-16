# coursebench-backend

course-bench 的后端服务

## Build

* Install Go
* Install Docker
* Install Docker-Compose

```shell
git clone git@github.com:ShanghaitechGeekPie/coursebench-backend.git
cd coursebench-backend
go build .
```

## Configure

```shell
cp config.json.example config.json
```

Edit ``config.json`` as you like.

Edit files in `build` directory as you like.

## Run

```shell
cd build
sudo docker-compose up -d
cd ..
./coursebench-backend
```