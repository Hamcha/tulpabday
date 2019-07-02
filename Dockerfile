FROM alpine:3.8 as kyotobuild

ENV LANG en_US.UTF-8
ENV KYOTOCABINET_VERSION 1.2.76

RUN apk update && apk upgrade
RUN apk --no-cache add libstdc++ lua go git
RUN apk --no-cache add --virtual build-dependencies build-base zlib-dev curl lua-dev
RUN curl -SLO http://fallabs.com/kyotocabinet/pkg/kyotocabinet-${KYOTOCABINET_VERSION}.tar.gz && \
    tar xzvf kyotocabinet-${KYOTOCABINET_VERSION}.tar.gz
RUN cd kyotocabinet-${KYOTOCABINET_VERSION} && \
    ./configure CFLAGS='-std=c++98' CXXFLAGS='-std=c++98' --enable-static && \
    make && \
    make install

ARG MODULE_PATH

ENV GOPROXY https://modules.fromouter.space
ENV GO111MODULE=on

WORKDIR /app

RUN apk add git; git clone https://github.com/hamcha/tulpabday; cd tulpabday; make

RUN apk --no-cache add lua

WORKDIR /app/tulpabday

CMD [ "/app/tulpabday/tulpaweb" ]