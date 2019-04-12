FROM eu.gcr.io/antha-images/golang:1.11-build

ADD . /go/src/github.com/antha-lang/antha
WORKDIR /go/src/github.com/antha-lang/antha
RUN mv .netrc $HOME/.netrc || true
RUN ./core-setup.sh
RUN set -ex && go get ./cmd/composer/ ./cmd/migrate/ ./cmd/elements/
RUN set -ex && go install ./cmd/composer/ ./cmd/migrate/ ./cmd/elements/
WORKDIR /app

# These are for the gitlab CI for elements:
ONBUILD ADD . /elements
ONBUILD ARG GIT_COMMIT_SHA
ONBUILD ENTRYPOINT /go/src/github.com/antha-lang/antha/cmd/elements/test.sh "$GIT_COMMIT_SHA"