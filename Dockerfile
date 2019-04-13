FROM eu.gcr.io/antha-images/golang:1.11-build

ARG GIT_COMMIT_SHA
ADD .netrc /
RUN mv /.netrc $HOME/.netrc || true
RUN mkdir /tmp/antha-core-build
WORKDIR /tmp/antha-core-build
RUN set -ex && go mod init antha-core-build && go get github.com/antha-lang/antha@$GIT_COMMIT_SHA
RUN set -ex && go install github.com/antha-lang/antha/cmd/...
ADD scripts/elements-test.sh /usr/local/bin/
WORKDIR /app

# These are for the gitlab CI for elements:
ONBUILD ADD . /elements
ONBUILD ARG GIT_COMMIT_SHA
ONBUILD ENTRYPOINT /go/src/github.com/antha-lang/antha/cmd/elements/test.sh "$GIT_COMMIT_SHA"