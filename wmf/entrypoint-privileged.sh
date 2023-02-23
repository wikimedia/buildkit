#!/bin/bash

set -eu -o pipefail

if [ -f /etc/buildkit/cni.json ]; then
    type=$(jq -r .type /etc/buildkit/cni.json)
    if [ "$type" == "bridge" ]; then
        name=$(jq -r .bridge /etc/buildkit/cni.json)

        # Ensure that Istio can't get its hand on outbound traffic
        # from build containers before masquerading. (T330433)
        echo Running iptables-legacy hack for https://phabricator.wikimedia.org/T330433
        /usr/sbin/iptables-legacy -t nat -I PREROUTING -i "$name" -j RETURN
    fi
fi

exec /usr/local/sbin/buildkitd \
  --oci-worker true \
  --oci-worker-gc \
  --containerd-worker false \
  "$@"
