FROM debian:buster-slim
RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY ./covidstats.json /app/creds.json
COPY ./bin/server /

ENV GOOGLE_APPLICATION_CREDENTIALS "/app/creds.json"

CMD /server
