FROM nginx:latest
WORKDIR /usr/src/app

COPY . .

COPY deploy/nginx.conf /etc/nginx/conf.d/default.conf

COPY src/ /usr/share/nginx/html


EXPOSE 80

