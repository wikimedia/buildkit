#!/bin/bash

set -eu -o pipefail

cd "$(dirname "$(dirname $0)")"
tmpdir="$(mktemp -d)"

cat Dockerfile wmf/Dockerfile > "$tmpdir"/Dockerfile.wmf
buildctl build \
  --frontend=dockerfile.v0 \
  --local context=. \
  --local dockerfile="$tmpdir" \
  --opt filename=Dockerfile.wmf \
  --opt target=wmf-production-privileged \
  "$@"
