FROM eu.gcr.io/antha-images/golang:1.12.4-build

ARG COMMIT_SHA
ARG NETRC
RUN printf "%s\n" "$NETRC" > $HOME/.netrc
RUN mkdir /antha
WORKDIR /antha
RUN set -ex && go get -u github.com/golangci/golangci-lint/cmd/golangci-lint   
RUN set -ex && go mod init antha && go mod edit "-require=github.com/antha-lang/antha@$COMMIT_SHA" && go mod download
RUN (set -ex && cd /go/pkg/mod/github.com/antha-lang/antha* && golangci-lint run --deadline=5m -E gosec -E gofmt ./...)
##RUN set -ex && golangci-lint run --deadline=5m -E gosec -E gofmt github.com/antha-lang/antha/...
RUN set -ex && go vet github.com/antha-lang/antha/...
RUN set -ex && go install github.com/antha-lang/antha/cmd/...
RUN set -ex && go test -c github.com/antha-lang/antha/cmd/elements
COPY scripts/. /antha/.
RUN ./antha-test.sh
RUN rm $HOME/.netr

# These are for the gitlab CI for elements:
ONBUILD ADD . /elements
