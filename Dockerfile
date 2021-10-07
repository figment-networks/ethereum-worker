# ------------------------------------------------------------------------------
# Builder Image
# ------------------------------------------------------------------------------
FROM golang AS build

WORKDIR /build

COPY ./go.mod .
COPY ./go.sum .

COPY .git .git
COPY ./Makefile ./Makefile
COPY ./api ./api
COPY ./client ./client
COPY ./cmd ./cmd
COPY ./structures ./structures
COPY ./transport ./transport

ENV GOARCH=amd64
ENV GOOS=linux

RUN make build-live

# ------------------------------------------------------------------------------
# Target Image
# ------------------------------------------------------------------------------

FROM alpine AS release

WORKDIR /app
COPY --from=build /build/ethereum-worker-live /app/ethereum-worker-live

RUN adduser --system --uid 1234 figment

RUN chmod a+x ./ethereum-worker-live
RUN chown -R figment /app/ethereum-worker-live

USER 1234

CMD ["./ethereum-worker-live"]
