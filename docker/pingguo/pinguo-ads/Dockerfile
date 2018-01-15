ENV MSF_ENV production-us


RUN set -xe \
       && echo "yac.enable_cli=1" >> /usr/local/etc/php/conf.d/docker-php-ext-yac.ini \

       # get code from s3
       && tar -xf $PROJECT_SRC \
       && find  -maxdepth 100 ! -name $PROJECT_DST|grep "/"|grep -E ".git"|xargs -i rm -rf {} \
       && rm -rf $PROJECT_SRC 
 


ADD ./php.ini /usr/local/etc/php/php.ini
