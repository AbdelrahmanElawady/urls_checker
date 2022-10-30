# URLs Checker

## How to run CLI

- Create a new config.toml file and insert your sites urls, for example:
  
```toml
[sites]
   [sites.threefold]
     url = "threefold.io"
   [sites.codescalers]
     url = "codescalers.com"
```

- build CLI `task build.cli`

- Run cmd `./bin/urls-checker-cli linkscheck --config=config.toml`

## How to run the app

### Run backend

```sh
task run.api
```

### Run frontend

```sh
cd fronend
npm install
npm run serve
```

## Run with docker

```sh
docker-compose up
```

## Run with helm

```sh
helm install urlschecker --debug ./chart 

export NODE_IP=$(kubectl get nodes --namespace default -o jsonpath="{.items[0].status.addresses[0].address}");
export NODE_PORT=$(kubectl get --namespace default -o jsonpath="{.spec.ports[0].nodePort}" services vue-serv);
echo http://$NODE_IP:$NODE_PORT

export NODE_IP=$(kubectl get nodes --namespace default -o jsonpath="{.items[0].status.addresses[0].address}");
export NODE_PORT=$(kubectl get --namespace default -o jsonpath="{.spec.ports[0].nodePort}" services server-serv);
echo http://$NODE_IP:$NODE_PORT
```

## Testing

Use this command to run the tests

```bash
task test
```
