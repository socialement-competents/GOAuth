export GOOS=linux
export GOARCH=amd64

unameOut="$(uname -s)"
case "${unameOut}" in
    CYGWIN*|MINGW*)     machine=Windows;;
    Linux*|Darwin*|*)   machine=Unix
esac

echo "building for $GOOS-$GOARCH on $machine"

rm -rf bin dist
mkdir bin dist

lambdas=""

while [ "$1" != "" ]
do
    lambdas="$lambdas $1"
    shift
done

if [ "$lambdas" == "" ]
then
    lambdas=$(ls lambdas)
fi

for lambda in $lambdas
do
    echo ""
    echo "building $lambda"
    go build -o bin/$lambda lambdas/$lambda/main.go || exit 1
    
    echo "zipping $lambda"
    if [ $machine = "Windows" ]; then
        # Grants the execute permission before zipping, on Windows
        # go get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip
        build-lambda-zip -o dist/$lambda.zip bin/$lambda
    else
        zip -j dist/$lambda.zip bin/$lambda
    fi

    # Add the /public subfolder to the .zip if it exists
    [ -e lambdas/$lambda/public ] && \
    echo "adding public/" && \
    cd lambdas/$lambda && \
    zip -ru ../../dist/$lambda.zip public && \
    cd ../..

    echo "built dist/$lambda.zip"
done
