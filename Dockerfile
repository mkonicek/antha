FROM eu.gcr.io/antha-images/golang:1.12.4-build

ARG COMMIT_SHA
ARG NETRC
RUN printf "%s\n" "$NETRC" > $HOME/.netrc
RUN mkdir /antha
WORKDIR /antha
RUN set -ex && go mod init antha && go mod edit "-require=github.com/antha-lang/antha@$COMMIT_SHA" && go mod download
# Do these builds to pre-warm the build cache. This makes a HUGE difference to performance in the cloud
RUN set -ex && go build github.com/antha-lang/antha/... github.com/Synthace/antha-runner/... github.com/Synthace/instruction-plugins/... && go install github.com/antha-lang/antha/cmd/...
RUN set -ex && go test -c github.com/antha-lang/antha/cmd/elements
COPY scripts/. /antha/.
RUN ./antha-test.sh
RUN rm $HOME/.netrc

# These are for the gitlab CI for elements:
ONBUILD ADD . /elements
