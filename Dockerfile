FROM alpine:3.20 AS runtime
ARG TARGETPLATFORM
RUN apk add --no-cache tini ca-certificates
WORKDIR /app

COPY ${TARGETPLATFORM}/app /usr/local/bin/app
COPY web/build /app/web
RUN chmod +x /usr/local/bin/app

ENV PORT=3000
ENV WEB_STATIC_PATH=/app/web

EXPOSE 3000

ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/app"]
