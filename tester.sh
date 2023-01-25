#/bin/sh
# This script will be used to test the server
url="http://localhost:8080/converter"
method="POST"
filepath="./data/data_lake/test.csv"
outputDir="./output/"
outputFile="response.txt"
# concatenate the output directory and the output file
outputPath=$outputDir$outputFile


# create an encoded b64 file
b64=$(base64 $filepath)

# make the request and save it as a txt file
curl -X $method $url -o $outputPath -s -d '{"filename":"test.csv", "content":"'"$b64"'", "extension":"csv", "uid": "test"}'

#curl -X $method $url -o $outputPath -s -d '{"filename":"test.csv", "content":$b64, "extension":"csv", "uid": "test"}'
# the -o is used to save the response to a file
# -s is used to hide the progress bar

# retrive the response from the txt file
response=$(cat $outputPath)

echo $response