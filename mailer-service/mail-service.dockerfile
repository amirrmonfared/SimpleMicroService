FROM alpine:3.18.3

RUN mkdir /app

COPY mailerApp /app

COPY templates /templates

CMD [ "/app/mailerApp" ]