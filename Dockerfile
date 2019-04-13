FROM eu.gcr.io/antha-images/golang:1.12-build

ARG COMMIT_SHA
ADD .netrc /
RUN mv /.netrc $HOME/.netrc || true
RUN mkdir /antha
WORKDIR /antha
RUN set -ex && go mod init antha && go get github.com/antha-lang/antha@$COMMIT_SHA && printf 'package main\nimport _ "github.com/antha-lang/antha/laboratory"\nfunc main() {}\n' > main.go && go build
RUN set -ex && go install github.com/antha-lang/antha/cmd/...
RUN set -ex && go test -c github.com/antha-lang/antha/cmd/elements
ADD scripts/elements-test.sh /antha/elements-test.sh
RUN rm $HOME/.netrc

# These are for the gitlab CI for elements:
ONBUILD ADD . /elements
ONBUILD ARG COMMIT_SHA
ONBUILD ENTRYPOINT /antha/elements-test.sh "$COMMIT_SHA"