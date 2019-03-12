#!/bin/bash
config=""

while getopts "c:" option 2>/dev/null
do
	case $option in
		c)
			config="$OPTARG"
			;;

		?)
			echo "[error] wrong arguments" 1>&2
            echo "Usage: start.sh [-c config_file_path]" 1>&2
            exit 1
            ;;
	esac
done

if [ -z $config ]
then
	echo "[error] wrong arguments" 1>&2
	echo "Usage: start.sh [-c config_file_path]" 1>&2
	exit 1
fi

go run $(ls | grep go | grep -v test | tr '\n' ' ') --config=$config
