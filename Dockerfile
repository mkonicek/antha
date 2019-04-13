FROM eu.gcr.io/antha-images/golang:1.11-build

ARG COMMIT_SHA
ADD .netrc /
RUN mv /.netrc $HOME/.netrc || true
RUN mkdir /antha
WORKDIR /antha
RUN set -ex && go mod init antha && go get github.com/antha-lang/antha@$COMMIT_SHA
RUN set -ex && go install github.com/antha-lang/antha/cmd/...
ADD scripts/elements-test.sh /antha/elements-test.sh
RUN rm $HOME/.netrc

# These are for the gitlab CI for elements:
ONBUILD ADD . /elements
ONBUILD ARG COMMIT_SHA
ONBUILD ENTRYPOINT /antha/elements-test.sh "$COMMIT_SHA"