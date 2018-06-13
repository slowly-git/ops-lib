#!/bin/bash

cd /home/worker/src/harbor

/bin/docker-compose stop

if [[ $? -eq 0 ]]
then
  /bin/docker run -it --name gc --rm --volumes-from registry vmware/registry:2.6.2-photon garbage-collect  /etc/registry/config.yml
  /bin/docker run --name gc --rm --volumes-from registry vmware/registry:2.6.2-photon garbage-collect  /etc/registry/config.yml
else
  echo "stop harbor faild, exit"
  exit 1
fi

/bin/docker-compose start 