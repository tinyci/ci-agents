from "golang:1.12"

run "mkdir -p /go/src/github.com/tinyci/ci-agents"
copy ".", "/go/src/github.com/tinyci/ci-agents", ignore_list: [".git", ".ca", ".logs", ".db", "build", "*.tar.gz"]
