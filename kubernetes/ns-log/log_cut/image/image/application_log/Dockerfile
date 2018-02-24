#FROM centos
FROM alpine 
MAINTAINER pinguoops <pinguo-ops@camera360.com>
RUN apk update && apk add dcron wget vim bash findutils coreutils rsync ca-certificates && rm -rf /var/cache/apk/*
ADD zoneinfo.tar.gz /
RUN mkdir -p /var/log/cron && mkdir -m 0644 -p /var/spool/cron/crontabs && touch /var/log/cron/cron.log && mkdir -m 0644 -p /etc/cron.d && rm -rf /etc/localtime && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
ADD root /etc/crontabs/
ADD clean_deployment_pods.sh /
ADD cut_app_pod.sh /
ADD application_log /
ADD application-entry.sh /
CMD ["/application-entry.sh"]
