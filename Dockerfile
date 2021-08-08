FROM debian:buster-slim
RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*


COPY ./server /
COPY ./static /static

CMD /server