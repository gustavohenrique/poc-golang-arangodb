# PoC ArangoDB and Golang

### Post

- [pt-BR Português](https://gustavohenrique.net/2019/08/banco-de-dados-em-grafo-com-arangodb-e-golang/)

### Fake data

```sh
virtualenv .venv
source .venv/bin/activate
pip install -r requirements.txt
python generate_fake_data.py
```

### ArangoDB

```sh
docker run -d --name arangodb -e ARANGO_ROOT_PASSWORD=root -v $PWD:/data --workdir=/data -p 8529:8529 arangodb
docker exec -it arangodb sh
sh> cd json && sh import_data.sh
```

### Webapp

```sh
GOOS=linux go build -o elearning -ldflags="-s -w" main.go
ARANGODB_HOST=127.0.0.1:8529 ./elearning
```

