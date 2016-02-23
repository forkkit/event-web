FROM alpine:3.2
ADD templates /templates
ADD event-web /event-web
WORKDIR /
ENTRYPOINT [ "/event-web" ]
