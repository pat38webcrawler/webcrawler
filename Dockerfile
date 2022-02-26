FROM debian:9


WORKDIR /app
ADD ./webcrawler /app

EXPOSE 8900

CMD ./webcrawler
