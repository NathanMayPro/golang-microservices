#/bin/sh

run_cmd_silent () {
    # echo "Running: ${1}"
    ${1} > /dev/null 2>&1
}

./launcher.sh & # launch the server in the background
./tester.sh # test the server

wait 3 # wait 3 seconds

run_cmd_silent "kill $(lsof -t -i:8080)"# kill the server the -s is used to kill the server in a silent way