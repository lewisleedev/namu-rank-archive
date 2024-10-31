FROM golang:1.23-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o namu-rank-archive

FROM alpine:latest
WORKDIR /app
COPY ./entrypoint.sh .
COPY --from=builder /app/namu-rank-archive .

ENV NAMU_RANK_DB=/data/ranks.db

RUN chmod +x /app/entrypoint.sh && \ 
    apk add --update tzdata && \
    mkdir -p /data

ENV TZ=Asia/Seoul

ENTRYPOINT ["./entrypoint.sh"]
