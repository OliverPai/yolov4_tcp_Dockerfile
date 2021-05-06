#!/bin/bash
DIR=$1

if [ ! -n "$DIR" ] ;then
    echo "you have not choice Application directory !"
    exit
fi

fswatch --event=Updated $DIR | while read file
do
  	mc mv ${file} myminio/yolov4tcp
	echo "${file} was updated"
done
