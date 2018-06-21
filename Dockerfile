FROM golang:latest

ARG SSH_DEP_PRIVATE_KEY

ENV UPX_VER 3.94
ENV APPDIR $GOPATH/src/github.com/l-vitaly/imhio-task
ENV BINNAME cfgm

RUN	env

# todo: check each installation for neediness
RUN apt-get update;\
    apt-get install -y xz-utils git autoconf automake libtool curl make g++ unzip && rm -rf /var/lib/apt/lists/*

# install UPX
ADD https://github.com/upx/upx/releases/download/v${UPX_VER}/upx-${UPX_VER}-amd64_linux.tar.xz /usr/local
RUN xz -d -c /usr/local/upx-${UPX_VER}-amd64_linux.tar.xz | \
    tar -xOf - upx-${UPX_VER}-amd64_linux/upx > /bin/upx && \
    chmod a+x /bin/upx

RUN git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"
RUN go get -u github.com/golang/dep/...

RUN mkdir ~/.ssh;\
    echo "$SSH_DEP_PRIVATE_KEY" > ~/.ssh/id_rsa;\
    chmod 600 ~/.ssh/id_rsa;\
    touch ~/.ssh/known_hosts;\
    ssh-keyscan gitlab.com >> ~/.ssh/known_hosts;

RUN mkdir -p ${APPDIR}

WORKDIR ${APPDIR}

ADD ./ .

RUN	make dep check build
RUN strip --strip-unneeded ${BINNAME}
RUN upx ${BINNAME}


FROM scratch

COPY --from=0 /go/src/github.com/l-vitaly/imhio-task/cfgm cfgm

EXPOSE 9000

CMD ["/cfgm"]
