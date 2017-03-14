FROM golang:alpine

ENV PATH /root/.yarn/bin:$PATH
RUN apk update \
  && apk add curl bash binutils tar git nodejs \
  && rm -rf /var/cache/apk/* \
  && /bin/bash \
  && touch ~/.bashrc \
  && curl -o- -L https://glide.sh/get | bash \
  && curl -o- -L https://yarnpkg.com/install.sh | bash \
  && apk del curl tar binutils

WORKDIR "${GOPATH}/src/github.com/duckclick/proxy"

ADD glide.yaml glide.yaml
ADD glide.lock glide.lock
RUN glide install

ADD frontend/package.json frontend/package.json
ADD frontend/yarn.lock frontend/yarn.lock
RUN cd frontend && yarn install

ADD . .
RUN go build -ldflags "-s -w"

ENV NODE_ENV production

CMD ["/bin/sh", "-c", "cd frontend && /root/.yarn/bin/yarn build && cd - && ./proxy"]
