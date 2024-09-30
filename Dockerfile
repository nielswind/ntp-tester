#FROM gcr.io/distroless/base-debian12 AS build-release-stage
FROM alpine AS build-release-stage

WORKDIR /app

COPY ntp_client_prometheus /app/ntp_client_prometheus

EXPOSE 2112

ENTRYPOINT ["/app/ntp_client_prometheus"]
