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
# 设置settings后无需指定-v参数（但settings更新时需重新打包）
COPY settings.yaml .

COPY --from=builder /build/main /app

EXPOSE 5000

CMD ["./main"]