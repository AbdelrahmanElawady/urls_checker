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
```

##### Project setup

```sh
npm install
```

##### Compiles and hot-reloads for development

```sh
npm run serve
```

## Testing

Use this command to run the tests

```bash
task test
```
