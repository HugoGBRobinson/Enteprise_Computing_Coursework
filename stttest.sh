#!/bin/sh
echo "{\"speech\":\"`base64 -i question1.wav`\"}" > input
JSON2=`curl -s -X POST -d @input localhost:3002/stt`
echo $JSON2

