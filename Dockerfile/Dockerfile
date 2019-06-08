FROM docker.io/busybox
MAINTAINER jikun.zhang <jikun.zhang>
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN chmod 755 /app/*
ENTRYPOINT ["/app/PrometheusAlert"]
