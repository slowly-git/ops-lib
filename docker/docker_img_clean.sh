#!/bin/bash

docker images | grep "weeks ago" |awk '{print $3}'|xargs -i docker rmi {}
