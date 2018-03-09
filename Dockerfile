FROM onsdigital/dp-docker-go

WORKDIR /app/

COPY ./build/dp-frontend-router .

ENTRYPOINT ./dp-frontend-router
