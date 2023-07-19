# 阶段一在golang镜像中编译go程序
FROM golang:1.19.3-alpine AS builder

LABEL stage=gobuilder

#更换阿里云镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache

#设置时区和环境变量
RUN apk add --no-cache tzdata

# 指定git构建分支
#ARG branch=master
# 安装git
#RUN apk add --no-cache git
#拉取代码
#ENV GOPROXY https://goproxy.cn,direct
#RUN cd src && git clone https://LeeKun42:Lijun0402@github.com/LeeKun42/mini-service-user.git

WORKDIR /go-build/src/

# 拷贝go项目代码到镜像中
COPY ./ /go-build/src

ENV GOPROXY https://goproxy.cn,direct

#编译go程序
RUN cd /go-build/src/ && go mod download && go build -ldflags="-s -w" -o /go/bin/mini-service-user main.go && cp env.yaml /go/bin/env.yaml

# 阶段二 在alpine镜像中运行编译好的go程序
FROM alpine

#更换阿里云镜像源
COPY --from=builder /etc/apk/repositories /etc/apk/repositories

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app/mini-service/

COPY --from=builder /go/bin/mini-service-user /app/mini-service/mini-service-user
COPY --from=builder /go/bin/env.yaml /app/mini-service/env.yaml

RUN touch stdout.log

CMD ["./mini-service-user"]

EXPOSE 8108
EXPOSE 8208

