#/bin/sh
# This script will be used to test the server
url="http://localhost:8080"
method="POST"
service="/save"
inputfileextension="pdf"
inputfilename="test-large-size"
inputfilepath="./data/data_lake/test-document/"$inputfilename"."$inputfileextension
savedirpath="./data/data_mart/test-document/"
outputDir="./output/"
# concatenate the output directory and the output file
urlWithQueryParams=$url$service"?path="$savedirpath
echo $urlWithQueryParams
# create an encoded b64 file
b64=$(base64 $inputfilepath)

# put in json format
echo '{"filename":"'"$inputfilename"'", "content":"'"$b64"'", "extension":"'"$inputfileextension"'", "uid": "test"}' > data.json

# first save file to data_mart
curl -X $method $urlWithQueryParams -d @data.json

rm data.json

echo 'file "'"$inputfilepath"'" saved in "'"$savedirpath$inputfilename.$inputfileextension"'"'

# then retrieve it from data_mart
method="GET"
service="/retrieve"
urlWithQueryParams=$url$service"?path=$savedirpath$inputfilename.$inputfileextension"

# # make the request and save it as a txt file
response=$(curl -X $method $urlWithQueryParams -o $outputDir$inputfilename'.'$inputfileextension) #-d '{"filename":"test", "content":"'"$b64"'", "extension":"pdf", "uid": "test"}'









#curl -X $method $url -o $outputPath -s -d '{"filename":"test.csv", "content":$b64, "extension":"csv", "uid": "test"}'
# the -o is used to save the response to a file
# -s is used to hide the progress bar

# retrive the response from the txt file
#response=$(cat $outputPath)

echo $response