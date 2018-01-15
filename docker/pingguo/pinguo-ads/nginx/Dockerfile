
FROM debian:latest

# -----------------------------------------------------------------------------
# Env src dir
# -----------------------------------------------------------------------------
ENV HOME /home/worker
ENV SRC_DIR $HOME/src
ENV nginx_version 1.12.2
ENV NGINX_INSTALL_DIR ${HOME}/nginx  
ENV LUAJIT_LIB ${HOME}/luajit/lib
ENV LUAJIT_INC ${HOME}/luajit/include/luajit-2.0

RUN mkdir -p ${SRC_DIR}

RUN set -xe \
        && extDeps=" \
                   libpcre3 \
                   curl \
                   unzip \
                   gcc \
                   make \
                   ca-certificates \
                   python \
                   python-pip \
                   " \
        && extDeps_dev=" \
                       python-setuptools \
                       libcurl4-openssl-dev \
                       libssl1.0-dev \ 
                       libpcre3-dev \
                       zlib1g-dev \                     
                       " \
        && apt-get update && apt-get install -y $extDeps $extDeps_dev --no-install-recommends && rm -rf /var/lib/apt/lists/* \
        && /usr/bin/pip install awscli \
    
        && cd ${SRC_DIR} \
        && curl -LJ -o nginx-${nginx_version}.tar.gz http://nginx.org/download/nginx-${nginx_version}.tar.gz \
        && curl -LJ -o nginx-http-concat.zip https://github.com/alibaba/nginx-http-concat/archive/master.zip \
        && curl -LJ -o nginx-logid.zip https://github.com/pinguo-liuzhaohui/nginx-logid/archive/master.zip \
        && curl -LJ -o ngx_devel_kit-0.3.0.tar.gz https://github.com/simpl/ngx_devel_kit/archive/v0.3.0.tar.gz \
        && curl -LJ -o lua-nginx-module-0.10.11.tar.gz https://github.com/openresty/lua-nginx-module/archive/v0.10.11.tar.gz \
        && curl -LJ -o LuaJIT-2.0.5.tar.gz http://luajit.org/download/LuaJIT-2.0.5.tar.gz \

        && tar zxf nginx-${nginx_version}.tar.gz \
        && unzip nginx-http-concat.zip -d nginx-http-concat \
        && unzip nginx-logid.zip -d nginx-logid \
        && tar zxf ngx_devel_kit-0.3.0.tar.gz \
        && tar zxf lua-nginx-module-0.10.11.tar.gz \
        && tar zxf LuaJIT-2.0.5.tar.gz \

        && cd LuaJIT-2.0.5 \
        && make PREFIX=${HOME}/LuaJIT-2.0.5 1>/dev/null \
        && make install PREFIX=${HOME}/LuaJIT-2.0.5 \

        && cd ${HOME} \
        && ln -s LuaJIT-2.0.5 luajit \

        && cd ${SRC_DIR}/nginx-${nginx_version} \
        && ./configure \
                --prefix=$NGINX_INSTALL_DIR \
                --with-http_stub_status_module \
                --with-http_ssl_module \
                --add-module=../nginx-http-concat/nginx-http-concat-master \
                --add-module=../nginx-logid/nginx-logid-master \
                --with-ld-opt="-Wl,-rpath,${HOME}/luajit/lib" \
                --add-module=../ngx_devel_kit-0.3.0 \
                --add-module=../lua-nginx-module-0.10.11 1>/dev/null \
        && make 1>/dev/null \
        && make install \

        && ln -s /home/worker/nginx/sbin/nginx /usr/local/bin/nginx \

        # deploy pinguo nginx config

        # deploy APP nginx conf


        # delete lib dev
        && rm -rf ${SRC_DIR}/nginx-* ${SRC_DIR}/ngx_devel_kit* ${SRC_DIR}/lua-nginx-module* ${SRC_DIR}/LuaJIT* \
        && apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false $extDeps_dev

EXPOSE 80

STOPSIGNAL SIGTERM

CMD ["nginx", "-g", "daemon off;"]
