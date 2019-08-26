# Estrutura de dados em Grafo usando ArangoDB com Golang

### Leia o artigo

[Texto sobre a PoC](POST.MD)

### Criando dados fake

Altere o script para aumentar a quantidade de objetos gerados.

```sh
virtualenv .venv
source .venv/bin/activate
pip install -r requirements.txt
python generate_fake_data.py
```

### Executando ArangoDB via Docker

```sh
docker run -d --name arangodb -e ARANGO_ROOT_PASSWORD=root -v $PWD:/data --workdir=/data -p 8529:8529 arangodb
docker exec -it arangodb sh
sh> cd json && sh import_data.sh
```

### Executando a Webapp

```sh
GOOS=linux go build -o elearning -ldflags="-s -w" main.go
ARANGODB_HOST=127.0.0.1:8529 ./elearning
```

