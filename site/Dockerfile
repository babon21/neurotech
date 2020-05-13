FROM nginx

COPY . /usr/share/nginx/html

RUN set -ex && \
    ln -sf /dev/stdout /var/log/nginx/access.log && \
    ln -sf /dev/stderr /var/log/nginx/error.log
