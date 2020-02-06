#!/bin/bash

getent passwd focal > /dev/null

STATUS=`echo $?`

if [ $STATUS -ne 0 ] ;
then
    useradd focal
fi