#!/bin/sh

mkdir -p /var/log/journal

/usr/local/bin/fluentd $@
