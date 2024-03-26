
FROM alpine:3.19

# Define the project name | 定义项目名称
ARG SERVICE_STYLE=bird
# Service name in snake format | 项目名称短杠格式
ARG SERVICE_DASH=bird
# The suffix after build or compile | 构建后缀
ARG PROJECT_BUILD_SUFFIX=rpc

ARG PROJECT=${SERVICE_STYLE}_${PROJECT_BUILD_SUFFIX}
# Define the config file name | 定义配置文件名
ARG CONFIG_FILE=${SERVICE_STYLE}.yaml

WORKDIR /app
ENV PROJECT=${PROJECT}
ENV CONFIG_FILE=${CONFIG_FILE}

COPY ./${PROJECT} ./
COPY ./etc ./etc

EXPOSE 9112

ENTRYPOINT ./${PROJECT}