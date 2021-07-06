## https://zhuanlan.zhihu.com/p/79949030s
FROM ubuntu:16.04
RUN apt-get update && apt-get install -y apt-utils \
  libjpeg-dev python-pip \
           && pip install --upgrade pip \
      && easy_install -U setuptools \
    && apt-get clean


## docker run -itd --name redis -p 6379:6379 redis:latest
## docker exec -it redis /bin/bash
## docker run -itd --name nsqlookupd -p 