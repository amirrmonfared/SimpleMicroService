FROM alpine:3.18.3

RUN mkdir /app

COPY loggerServiceApp /app

CMD [ "/app/loggerServiceApp" ]