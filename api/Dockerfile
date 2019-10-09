FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go build -o main .

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
COPY --from=builder /build/blog.config.toml /app/
COPY --from=builder /build/migrations/* /app/migrations/
WORKDIR /app
CMD ["./main"]