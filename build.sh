export GOOS=linux
export GOARCH=amd64

unameOut="$(uname -s)"
case "${unameOut}" in
    CYGWIN*|MINGW*)     machine=Windows;;
    Linux*|Darwin*|*)   machine=Unix
esac

echo "building for $GOOS-$GOARCH"

go build -o bin/hellolambda lambdas/hellolambda/main.go

echo "zipping on $machine"

if [ $machine = "Windows" ]; then
    build-lambda-zip -o dist/hellolambda.zip bin/hellolambda
else
    zip -j dist/hellolambda.zip bin/hellolambda
fi
