# CDK Lambda GO

- CDK App that deploy golang Lambda Function 

## Useful commands

### build

```zsh
# lambda functio build
cd lambda; GOOS=linux go build -o bin/main

# project build
npm run build
```

### Lambda Function Local invoke

```zsh
cdk synth --no-staging > template.yaml
sam local invoke MyFuncHandlerXXX -e lambda/sample-event.json
```