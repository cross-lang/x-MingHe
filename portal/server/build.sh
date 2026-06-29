#!/bin/bash

X_NAME="MingHe_portal"
IMAGE_NAME="MingHe_portal:latest"

ENV=$1

if [ -z "$ENV" ]; then
  echo "❌ 请输入环境参数: dev | test | prod"
  echo "示例: ./run.sh dev"
  exit 1
fi

# ======================
# 构建镜像
# ======================
echo "📦 正在构建镜像: ${IMAGE_NAME}"
if ! docker build -t "${IMAGE_NAME}" .; then
  echo "❌ 镜像构建失败，停止执行！"
  exit 1
fi

# ======================
# 清理旧容器
# ======================
echo "🛑 停止旧容器（如果存在）：${X_NAME}"
docker stop ${X_NAME} 2>/dev/null
echo "🛑 删除旧容器（如果存在）：${X_NAME}"
docker rm ${X_NAME} 2>/dev/null

# ======================
# 根据环境执行不同逻辑
# ======================
case "$ENV" in
  dev)
    echo "🚀 启动 DEV 环境"
    docker run -d --name ${X_NAME} -v ./log:/MingHe/portal/log -v ./config_dev.yaml:/MingHe/portal/config.yaml -p 8088:8088 ${IMAGE_NAME}
    ;;

  test)
    echo "🚀 启动 TEST 环境"
    docker run -d --name ${X_NAME} -v ./log:/MingHe/portal/log -v ./config_test.yaml:/MingHe/portal/config.yaml -p 8088:8088 ${IMAGE_NAME}
    ;;

  prod)
    echo "🚀 启动 PROD 环境"
    docker run -d --name ${X_NAME} -v ./log:/MingHe/portal/log -v ./config_prod.yaml:/MingHe/portal/config.yaml -p 8088:8088 ${IMAGE_NAME}
    ;;

  *)
    echo "❌ 不支持的环境: $ENV"
    echo "可选值: dev | test | prod"
    exit 1
    ;;
esac

# ======================
# 清理镜像
# ======================
docker image prune -af
