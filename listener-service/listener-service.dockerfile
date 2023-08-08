FROM alpine:3.18.3

RUN mkdir /app

COPY listenerApp /app

CMD [ "/app/listenerApp"]