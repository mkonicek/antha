FROM eu.gcr.io/antha-images/golang:1.12.4-build AS build

ARG COMMIT_SHA
ARG NETRC
RUN printf "%s\n" "$NETRC" > $HOME/.netrc
WORKDIR /antha
RUN set -ex && go mod init antha && go mod edit "-require=github.com/antha-lang/antha@$COMMIT_SHA" && go mod download
RUN set -ex && go install github.com/antha-lang/antha/cmd/...
RUN set -ex && go test -c github.com/antha-lang/antha/cmd/elements
COPY scripts/. /antha/.
RUN ./antha-test.sh

# Final stage: drop the layers with the netrc credentials. we still need the go
# compiler, etc, so use the same base image.
FROM eu.gcr.io/antha-images/golang:1.12.4-build 
WORKDIR /antha
COPY --from=build /antha/. .
COPY --from=build /go/. /go/.

# These are for the gitlab CI for elements:
ONBUILD ADD . /elements
