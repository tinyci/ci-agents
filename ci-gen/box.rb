GO_VERSION = "1.13"
SWAGGER_VERSION = "v0.18.0"
PROTOC_VERSION = "3.7.1"
TIMEZONE = "Etc/UTC"

EXTRA_PACKAGES = %w[
  curl 
  wget
  gnupg 
  git 
  mercurial 
  build-essential
  sudo
  nginx
  libnss3-tools
  unzip
]

from "ubuntu:20.04"

env TZ: TIMEZONE
run "ln -sf /usr/share/zoneinfo/#{TIMEZONE} /etc/localtime"

after do
  if getenv("PACKAGE_FOR_CI") != ""
    run "apt-get clean"
    run "rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /root/.cache"
    flatten
  end
end

run %Q[perl -i.bak -pe 's!//(security|archive).ubuntu.com!//#{getenv("APT_MIRROR").length > 0 ? getenv("APT_MIRROR") : "mirror.pnl.gov"}!g' /etc/apt/sources.list]

run "apt-get update && apt-get dist-upgrade -y && apt-get install #{EXTRA_PACKAGES.join(" ")} -y"

env GOPATH: "/go",
    PATH: %w[
      /go/bin
      /usr/local/go/bin
      /usr/local/sbin
      /usr/local/bin
      /usr/sbin
      /usr/bin
      /sbin
      /bin
    ].join(":"),
    TINYCI_CONFIG: "./.config",
    TZ: "Etc/UTC",
    TESTING: getenv("TESTING"),
    CAROOT: "/var/ca"

run "curl -sSL https://dl.google.com/go/go#{GO_VERSION}.linux-amd64.tar.gz | tar xz -C /usr/local"
run "mkdir /go"

run "go get -u github.com/golang/protobuf/protoc-gen-go"
run "go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc"

protoc_fn = "protoc-#{PROTOC_VERSION}-linux-x86_64.zip"

run "wget https://github.com/protocolbuffers/protobuf/releases/download/v#{PROTOC_VERSION}/#{protoc_fn}"
run "unzip '#{protoc_fn}' -d /usr && rm -f '#{protoc_fn}'"
run "curl -sSL 'https://github.com/go-swagger/go-swagger/releases/download/#{SWAGGER_VERSION}/swagger_linux_amd64' >/go/bin/swagger && chmod +x /go/bin/swagger"
