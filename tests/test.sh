#!/bin/sh

function kill_fwego {
    kill `ps ax | grep fwego | grep -v grep | awk {'print $1'}` || panic
}

function panic {
    echo "Fail"
    kill_fwego
    exit 1
}

# first run application
./fwego -root-path . &

# second. build test app
pushd tests > /dev/null 2>&1
go build test_fwego.go || panic
./test_fwego || panic
popd > /dev/null 2>&1

# third. kill fwego
kill_fwego

exit 0
