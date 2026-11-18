FROM alpine:3.22
ARG TARGETPLATFORM

RUN apk add --no-cache tini ca-certificates tzdata \
	&& addgroup -S -g 10001 app \
	&& adduser -S -D -H -u 10001 -G app app \
	&& mkdir -p /data \
	&& chown app:app /data

WORKDIR /app
COPY --chmod=0755 ${TARGETPLATFORM}/app /usr/local/bin/app
COPY web/build /app/web

ENV PORT=3000 \
	DATA_DIR=/data \
	WEB_STATIC_PATH=/app/web

USER app:app
VOLUME /data
EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 \
	CMD wget -q -O /dev/null "http://127.0.0.1:${PORT}/" || exit 1

ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/app"]
