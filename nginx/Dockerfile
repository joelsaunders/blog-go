FROM nginx:1.11.9-alpine
EXPOSE 80
RUN rm -fr /var/www
RUN mkdir -p /var/www
WORKDIR /etc/nginx
RUN rm nginx.conf
COPY ./ /etc/nginx/
RUN mv www/* /var/www
RUN rm -fr www
CMD ["nginx"]
