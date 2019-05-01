FROM eu.gcr.io/antha-images/golang:1.12.4-build AS build

ARG NETRC
RUN printf "%s\n" "$NETRC" > $HOME/.netrc
WORKDIR /antha
# layer for caching build dependencies
COPY go.mod go.sum ./
RUN go mod edit -dropreplace=github.com/Synthace/antha-runner -dropreplace=github.com/Synthace/instruction-plugins
RUN go mod download

# Now build antha commands
COPY . .
# repeat since we copied these again
RUN go mod edit -dropreplace=github.com/Synthace/antha-runner -dropreplace=github.com/Synthace/instruction-plugins
RUN go install ./cmd/...  && \
    go test -c ./cmd/elements
RUN scripts/antha-test.sh

# Final stage: drop the layers with the netrc credentials. we still need the go
# compiler, etc, so use the same base image.
FROM eu.gcr.io/antha-images/golang:1.12.4-build 
WORKDIR /antha
COPY --from=build /antha/. .
COPY --from=build $HOME/.cache  $HOME/
COPY --from=build /go/. /go/.

# These are for the gitlab CI for elements:
ONBUILD ADD . /elements
