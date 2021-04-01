#!/bin/bash

set -e

WORKSPACE=${GOPATH}/src/dinit
BUILD_PATH=${WORKSPACE}/build

cd $WORKSPACE/cmd
DIRS=$(find * -maxdepth 0 -type d)
pushd $(pwd) > /dev/null
for app_dir in $DIRS;do
    FILES=$(find $app_dir -name 'Makefile')
    for app_makefile in $FILES;do
        app_makefile_path=$(pwd)/$app_makefile
        if [ -f $app_makefile_path ];then
            pushd $(pwd) > /dev/null
            cd $(dirname $app_makefile_path)
            make
            if [ $? -ne 0 ];then
                exit
            fi
            popd > /dev/null
        fi
    done
done
popd > /dev/null
