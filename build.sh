GOOS=linux
GOARCH=amd64

echo "building for $GOOS-$GOARCH"

go build -o dist/createschema/main lambdas/createschema/main.go
zip dist/createschema/main.zip dist/createschema/main

go build -o dist/hellolambda/main lambdas/hellolambda/main.go
zip dist/hellolambda/main.zip dist/hellolambda/main
