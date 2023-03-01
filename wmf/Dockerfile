# This Dockerfile is meant to be concatenated with upstream's and then
# executed in the top-level directory context.

# There are two targets defined:
# * wmf-production-rootless
#   Note that in rootless mode, there will be less isolation between the build
#   container processes that the runc worker spawns and the network namespace
#   they in live under. There is also the possibility of network space
#   collisions (port allocation, etc.).
# 
# * wmf-production-privileged
#   In rootfull mode, a privileged buildkitd container is required for network
#   namespace creation and cgroup management.
#
# Example for rootless mode execution:
#   docker run --security-opt seccomp=unconfined --security-opt apparmor=unconfined
#
# Example for root mode:
#   docker run --privileged

# See https://github.com/moby/buildkit/blob/master/docs/rootless.md
# See https://phabricator.wikimedia.org/T307810 for a brief evaluation of the
# attack surface from build processes in rootless mode.


FROM docker-registry.wikimedia.org/bullseye:20221002 AS wmf-production-base

ARG PACKAGES="ca-certificates wmf-certificates git libcap2-bin runc uidmap iptables"
ARG DEBUG_PACKAGES=""
#ARG DEBUG_PACKAGES="iproute2 strace procps net-tools curl"

RUN --mount=type=cache,target=/var/cache/apt \
    --mount=type=cache,target=/var/lib/apt/lists \
    apt-get update && \
    apt-get install -yqq --no-install-recommends $PACKAGES $DEBUG_PACKAGES

# This buildkitd image is intended to run in rootless mode using rootlesskit.
#
# Note that in rootless mode, there will be less isolation between the build
# container processes that the runc worker spawns and the network namespace
# they in live under. There is also the possibility of network space
# collisions (port allocation, etc.).
#
# Example for rootless mode execution:
#   docker run --security-opt seccomp=unconfined --security-opt apparmor=unconfined
#
# See https://github.com/moby/buildkit/blob/master/docs/rootless.md
# See https://phabricator.wikimedia.org/T307810 for a brief evaluation of the
# attack surface from build processes in rootless mode.
FROM wmf-production-base AS wmf-production-rootless

ARG BUILDKIT_HOME=/var/lib/buildkit
ARG BUILDKIT_UID=1000

COPY --from=rootless /usr/bin/buildkitd /usr/bin/buildctl /usr/local/sbin/
COPY --from=rootless /usr/bin/rootlesskit /usr/local/sbin/
COPY wmf/entrypoint-rootless.sh /usr/local/bin/entrypoint.sh

# Much of this section was adapted from https://github.com/moby/buildkit/blob/v0.10/Dockerfile
RUN groupadd -r -g $BUILDKIT_UID buildkit \
    && useradd -r -u $BUILDKIT_UID -g $BUILDKIT_UID -d $BUILDKIT_HOME buildkit \
    && mkdir -p /run/user/$BUILDKIT_UID $BUILDKIT_HOME/.local/tmp $BUILDKIT_HOME/.local/share/buildkit \
    && chown -R buildkit:buildkit /run/user/$BUILDKIT_UID $BUILDKIT_HOME \
    && echo buildkit:100000:65536 | tee /etc/subuid | tee /etc/subgid

# Remove suid from newuidmap and newgidmap commands and add
# cap_setuid/cap_setgid. The latter is needed to allow writing to uid_map and
# gid_map and for some strange reason suid seems to interfere. See
# https://man7.org/linux/man-pages/man7/user_namespaces.7.html
RUN chmod u-s /usr/bin/newuidmap /usr/bin/newgidmap \
    && setcap cap_setuid=ep /usr/bin/newuidmap \
    && setcap cap_setgid=ep /usr/bin/newgidmap

USER $BUILDKIT_UID:$BUILDKIT_UID
ENV HOME $BUILDKIT_HOME
ENV USER buildkit
ENV XDG_RUNTIME_DIR=/run/user/$BUILDKIT_UID
ENV TMPDIR=$BUILDKIT_HOME/.local/tmp

ENTRYPOINT [ "/usr/local/bin/entrypoint.sh" ]

# This buildkitd image requires a privileged container but can offer better
# process isolation and network isolation (via CNI) between runc
# worker-spawned processes.
FROM wmf-production-base AS wmf-production-privileged

COPY --from=buildkit-linux /usr/bin/buildkitd /usr/bin/buildctl /usr/local/sbin/
COPY wmf/entrypoint-privileged.sh /usr/local/bin/entrypoint.sh

# Copy in CNI binaries from cni-plugins stage in upstream Dockerfile
COPY --link --from=cni-plugins /opt/cni/bin /opt/cni/bin
COPY hack/fixtures/cni.json /etc/buildkit/cni.json

ENTRYPOINT [ "/usr/local/bin/entrypoint.sh" ]
