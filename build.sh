#!/bin/bash

# 当前目录下editor文件夹下所有js文件sed -i 's/127.0.0.1:9090/127.0.0.1:9099/g' 
find editor -name "*.js" -type f -exec sed -i '' 's/127.0.0.1:9090/127.0.0.1:9099/g' {} \;

docker compose up -d