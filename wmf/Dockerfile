# This Dockerfile is meant to be concatenated with upstream's and then
# executed in the top-level directory context

# This buildkitd image is intended for rootless mode execution using
# rootlesskit. It must be run via:
#   docker run --security-opt seccomp=unconfined --security-opt apparmor=unconfined
# See https://github.com/moby/buildkit/blob/master/docs/rootless.md
# See https://phabricator.wikimedia.org/T307810 for a brief evaluation of the
# attack surface from build processes
FROM docker-registry.wikimedia.org/bullseye:20221002 AS wmf-production

ARG BUILDKIT_HOME=/var/lib/buildkit
ARG BUILDKIT_UID=1000

ARG PACKAGES="ca-certificates wmf-certificates git libcap2-bin runc uidmap"

RUN --mount=type=cache,target=/var/cache/apt \
    --mount=type=cache,target=/var/lib/apt/lists \
    apt-get update && \
    apt-get install -yqq --no-install-recommends $PACKAGES

COPY --from=rootless /usr/bin/buildkitd /usr/bin/buildctl /usr/local/sbin/
COPY --from=rootless /usr/bin/rootlesskit /usr/local/sbin/
COPY wmf/entrypoint.sh /usr/local/bin/

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
