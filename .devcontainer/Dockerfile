FROM golang:1.19beta1-bullseye

RUN \
  # install golangci-lint
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
  # install godoc
  && go install golang.org/x/tools/cmd/godoc@latest \
  # install dependencies for VSCode golang extension
  # cf. https://github.com/golang/vscode-go/blob/master/docs/tools.md
  && go install github.com/cweill/gotests/gotests@latest \
  && go install github.com/fatih/gomodifytags@latest \
  && go install github.com/josharian/impl@latest \
  && go install github.com/haya14busa/goplay/cmd/goplay@latest \
  && go install honnef.co/go/tools/cmd/staticcheck@latest \
  && go install golang.org/x/tools/gopls@latest \
  && go install github.com/ramya-rao-a/go-outline@latest \
  && go install github.com/go-delve/delve/cmd/dlv@latest
