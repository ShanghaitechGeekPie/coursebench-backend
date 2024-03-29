# coursebench-backend

The backend service of [GeekPie_CourseBench](https://coursebench.geekpie.club/).

## Build

* Install Go
* Install Docker
* Install Docker-Compose

```shell
git clone git@github.com:ShanghaitechGeekPie/coursebench-backend.git
cd coursebench-backend
mkdir bin
go build -o bin/coursebench-backend ./cmd/coursebench-backend/
go build -o bin/cmd_tools ./cmd/cmd_tools/
```

If you are using Windows, then the two `build` instruction should be

```shell
go build -o bin/coursebench-backend.exe ./cmd/coursebench-backend/
go build -o bin/cmd_tools.exe ./cmd/cmd_tools/
```



## Configure

```shell
cp config.json.example config.json
```

Edit ``config.json`` as you like.

Edit files in `build` directory as you like.

## Import data (optional)

```shell
bin/cmd_tools import_course <course data path>
```

## Run

```shell
cd build
docker-compose up -d
cd ..
./bin/coursebench-backend
```

If you are using Windows, the last instruction should be

```shell
./bin/coursebench-backend.exe
```

