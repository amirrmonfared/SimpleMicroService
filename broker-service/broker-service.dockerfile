FROM alpine:3.18.3

RUN mkdir /app

COPY brokerApp /app

CMD [ "/app/brokerApp" ]