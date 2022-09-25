# ===============================================================================
# Global arguments
# ===============================================================================
ARG project_name=trading-service

# ===============================================================================
# Stage 1: Compile the executable files from source
# ===============================================================================
FROM golang:1.19.1-alpine3.16 AS build-env

ARG project_name
ARG app_name

COPY . /srv/$project_name/src

WORKDIR /srv/$project_name

# Compile the go executable file
RUN apk add build-base \
    && mkdir -p /srv/$project_name/build \
    && cd /srv/$project_name/src \
    && go build -o /srv/$project_name/build/$app_name -tags musl ./cmd/$app_name/*

# ===============================================================================
# Stage 2: Build the runtime image
# ===============================================================================
FROM alpine:3

ARG project_name
ARG app_name

# Copy the compiled go executable files from build-env
WORKDIR /srv/$project_name

COPY --from=build-env /srv/$project_name/build /srv/$project_name/bin
COPY --from=build-env /srv/$project_name/src/docker-entrypoint.sh /usr/local/bin/

EXPOSE 8080

# Volume config files
VOLUME ["/srv/$project_name/config/config.yaml"]

# Set the environment variables for docker-entrypoint.sh
ENV PROJECT_NAME=$project_name
ENV APP_NAME=$app_name

ENTRYPOINT ["docker-entrypoint.sh"]
