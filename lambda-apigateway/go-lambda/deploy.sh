version=$1
if [ -z "$version" ]; then
    echo "Please provide a version argument."
    exit 1
fi


aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 383660184915.dkr.ecr.us-west-2.amazonaws.com

docker build -t lambda-gateway:$version .

docker tag lambda-gateway:$version 383660184915.dkr.ecr.us-west-2.amazonaws.com/lambda-gateway:$version

docker push 383660184915.dkr.ecr.us-west-2.amazonaws.com/lambda-gateway:$version


aws lambda update-function-code --region us-west-2  --function-name lambda-go \
  --image-uri 383660184915.dkr.ecr.us-west-2.amazonaws.com/lambda-gateway:$version
