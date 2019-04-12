FROM eu.gcr.io/antha-images/golang:1.11-build

ADD . /go/src/github.com/antha-lang/antha
WORKDIR /go/src/github.com/antha-lang/antha
RUN mv .netrc $HOME/.netrc || true
RUN set -ex && go get ./cmd/composer/ ./cmd/migrate/ ./cmd/elements/
RUN set -ex && go install ./cmd/composer/ ./cmd/migrate/ ./cmd/elements/
WORKDIR /app
