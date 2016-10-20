FROM ubuntu:16.04

WORKDIR /app/

COPY ./build/dp-frontend-router .

ENTRYPOINT ./dp-frontend-router
