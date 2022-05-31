FROM gcr.io/distroless/static
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY ./connector-spacelift-${TARGETOS}-${TARGETARCH} /app/connector-spacelift
ENTRYPOINT [ "/app/connector-spacelift" ]