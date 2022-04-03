# create by xiexianbin, Github Action for go-gcrmirrors
FROM alpine:latest

# Dockerfile build cache
ENV REFRESHED_AT 2022-04-02

LABEL "com.github.actions.name"="go-gcrmirrors"
LABEL "com.github.actions.description"="Github Action for auto generate https://github.com/kbcx/mirrors.kb.cx api json from https://github.com/kbcx/gcr.io."
LABEL "com.github.actions.icon"="home"
LABEL "com.github.actions.color"="green"
LABEL "repository"="http://github.com/x-actions/go-gcrmirrors"
LABEL "homepage"="http://github.com/x-actions/go-gcrmirrors"
LABEL "maintainer"="xiexianbin<me@xiexianbin.cn>"

LABEL "Name"="Github Action for go-gcrmirrors"
LABEL "Version"="v1.0.0"

ENV LC_ALL C.UTF-8
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US.UTF-8

RUN apk update && apk add --no-cache git git-lfs bash wget curl openssh-client tree && rm -rf /var/cache/apk/*

RUN mkdir /usr/local/go-gcrmirrors/ && \
    cd /usr/local/go-gcrmirrors/ && \
    curl -s https://api.github.com/repos/x-actions/go-gcrmirrors/releases/latest | \
    sed -r -n '/browser_download_url/{/linux.tar.gz/{s@[^:]*:[[:space:]]*"([^"]*)".*@\1@g;p;q}}' | xargs wget && \
    tar xzf *linux.tar.gz -C /usr/local/go-gcrmirrors/ && \
    cp /usr/local/go-gcrmirrors/gcrmirrors_*_linux/gcrmirrors /usr/local/bin/ && \
    rm -rf /usr/local/go-gcrmirrors/

ADD entrypoint.sh /
RUN chmod +x /entrypoint.sh

WORKDIR /github/workspace
ENTRYPOINT ["/entrypoint.sh"]
