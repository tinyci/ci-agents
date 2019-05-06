#!/bin/sh

set -e

package="github.com/tinyci/ci-agents"

dirs=$(go list ./... | sed -e "s!${package}!.!g" | grep -vE '^.$' | grep -v ./vendor)
files=$(find . -type f -name '*.go' 2>/dev/null | grep -vE '^.$' | grep -v ./vendor)

echo "Testing packages:"
for dir in ${dirs}
do
  echo ${dir} | sed -e "s!^./!${package}/!g"
done

echo
echo "-------------------- STARTING CHECKS --------------------"
echo

echo "Running gofmt..."
set +e
out=$(gofmt -s -l ${files})
set -e
# gofmt can include empty lines in its output
if [ "`echo \"${out}\" | sed '/^$/d' | wc -l`" -gt 0 ]
then
  echo 1>&2 "gofmt errors in:"
  echo 1>&2 "${out}"
  exit 1
fi

echo "Running ineffassign..."
[ -n "`which ineffassign`" ] || go get github.com/gordonklaus/ineffassign &>/dev/null
for i in ${dirs}
do
  ineffassign $i
done

echo "Running golint..."
[ -n "`which golint`" ] || go get golang.org/x/lint/golint &>/dev/null
set +e
out=$(golint ./... | grep -vE '^vendor')
set -e
if [ "`echo \"${out}\" | sed '/^$/d' | wc -l`" -gt 0 ]
then
  echo 1>&2 "golint errors in:"
  echo 1>&2 "${out}"
  exit 1
fi

echo "Running govet..."
set +e
out=$(go vet -composites=false ${dirs} 2>&1 | grep -v vendor)
set -e

if [ "`echo \"${out}\" | sed '/^$/d' | wc -l`" -gt 0 ]
then
  echo 1>&2 "go vet errors in:"
  echo 1>&2 "${out}"
  exit 1
fi

echo "Running gocyclo..."
[ -n "`which gocyclo`" ] || go get github.com/fzipp/gocyclo &>/dev/null
set +e
out=$(gocyclo -over 15 . | grep -v vendor)
set -e
if [ "`echo \"${out}\" | sed '/^$/d' | wc -l`" -gt 0 ]
then
  echo 1>&2 "gocycle errors in:"
  echo 1>&2 "${out}"
  exit 1
fi

echo "Running misspell..."
[ -n "`which misspell`" ] || go get github.com/client9/misspell/... &>/dev/null
set +e
out=$(misspell -locale US -error -i exportfs ${dirs} | grep -vE '^vendor')
set -e
if [ "`echo \"${out}\" | sed '/^$/d' | wc -l`" -gt 0 ]
then
  echo 1>&2 "misspell errors in:"
  echo 1>&2 "${out}"
  exit 1
fi
