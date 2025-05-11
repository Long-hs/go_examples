FROM golang:alpine AS builder

# 构建可执行文件
ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

WORKDIR /build
ADD go.mod .
ADD go.sum .
ADD main.go .
RUN go build -o main


FROM scratch
WORKDIR /app
COPY --from=builder /build/main /app
CMD ["./main"]