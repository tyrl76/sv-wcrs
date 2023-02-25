FROM ubuntu:22.04

RUN apt-get update; apt-get install -y ca-certificates;

COPY ./dist/ /app/

RUN cd /app/
RUN mkdir /app/config
RUN chmod 777 -R /app/

ENTRYPOINT ["/app/sn-wcrs"]