GO_VERSION = "1.13"
POSTGRES_VERSION = "11"
SWAGGER_VERSION = "v0.18.0"
PROTOC_VERSION = "3.7.1"

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

from "ubuntu:19.04"

run %Q[perl -i.bak -pe 's!//(security|archive).ubuntu.com!//#{getenv("APT_MIRROR").length > 0 ? getenv("APT_MIRROR") : "mirror.pnl.gov"}!g' /etc/apt/sources.list]

run "apt-get update && apt-get dist-upgrade -y && apt-get install #{EXTRA_PACKAGES.join(" ")} -y"

run "curl -sSL https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -"
run "echo 'deb http://apt.postgresql.org/pub/repos/apt/ disco-pgdg main' | tee -a /etc/apt/sources.list.d/postgresql.list"

run "ln -s /usr/share/zoneinfo/Etc/UTC /etc/localtime"

env GOPATH: "/go",
    PATH: %W[
      /usr/lib/postgresql/#{POSTGRES_VERSION}/bin
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

run "apt-get update -qq && apt-get install postgresql-#{POSTGRES_VERSION} -y -qq"

run "curl -sSL https://dl.google.com/go/go#{GO_VERSION}.linux-amd64.tar.gz | tar xz -C /usr/local"
run "mkdir /go"

run "go get github.com/erikh/migrator"
run "go get github.com/FiloSottile/mkcert"
run "go get -u github.com/golang/protobuf/protoc-gen-go"
run "go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc"

protoc_fn = "protoc-#{PROTOC_VERSION}-linux-x86_64.zip"

run "wget https://github.com/protocolbuffers/protobuf/releases/download/v#{PROTOC_VERSION}/#{protoc_fn}"
run "unzip '#{protoc_fn}' -d /usr"
run "curl -sSL 'https://github.com/go-swagger/go-swagger/releases/download/#{SWAGGER_VERSION}/swagger_linux_amd64' >/go/bin/swagger && chmod +x /go/bin/swagger"

if !$imported
  copy '.config/nginx.conf', '/etc/nginx/nginx.conf'
  copy 'entrypoint.sh', '/'
  run "chmod 755 /entrypoint.sh"
end

if !$imported
  entrypoint '/entrypoint.sh'
end
