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

## Testing

Use this command to run the tests

```bash
go test -v ./...
```
