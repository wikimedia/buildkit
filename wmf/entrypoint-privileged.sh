#!/bin/bash

exec /usr/local/sbin/buildkitd \
  --oci-worker true \
  --oci-worker-gc \
  --containerd-worker false \
  "$@"
