# 使用官方的 Redis 镜像
FROM redis:latest

# 设置工作目录
WORKDIR /data

# CMD ["redis-server"]
CMD ["redis-server"]


