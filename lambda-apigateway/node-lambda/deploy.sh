zip -r function.zip .

aws lambda update-function-code  --region us-west-2 --function-name lambda-autorizer --zip-file fileb://function.zip