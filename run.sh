#!/bin/bash

if [[ "$1" == "" ]]; then
	echo "usage: ./run.sh console_file_name"
	exit 1
fi

go=`which go`

root_dir=$(cd $(dirname "$0") && pwd)
console_dir=$root_dir/console

source_file=$console_dir/"$1".go

if [ ! -f "$source_file" ]; then
	echo "error: $source_file not exists"
	exit 2
fi

command="/tmp/$1"
go build -o $command $source_file
if [ ! -x "$command" ]; then
	echo "error: build command error"
	exit 3
fi

cd $root_dir
($command ${@:2})