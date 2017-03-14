# Proxy

## Running

Compile the assets:
```sh
cd frontend/
yarn install
yarn build
```

Run the server:
```sh
go run main.go
```

## Development

__Proxy__ uses `glide` for dependency management, to setup the project do the following:

```sh
brew install glide
```

__Note__: Read here (https://github.com/Masterminds/glide) for a different OS

```sh
glide install
```

To add new dependecies do:

```sh
glide get <dependency>
```

Build an image with docker

```sh
docker build -t proxy .
```

Run a docker container

```sh
docker run --rm -t -p 7275:7275 -e "BEAK_HOST=http://localhost:7276" proxy
```

## Examples

```sh
curl -v "http://localhost:7275/__duckclick__/configure" -d '{
  "url": "http://todomvc.com",
  "host": "todomvc.com",
  "current_path": "examples/react"
}'

curl -v "http://localhost:7275/node_modules/todomvc-app-css/index.css" --cookie "duckclick.proxy.configure=eyJ1cmwiOiJodHRwOi8vdG9kb212Yy5jb20iLCJob3N0IjoidG9kb212Yy5jb20iLCJjdXJyZW50X3BhdGgiOiJleGFtcGxlcy9yZWFjdCJ9"
```
