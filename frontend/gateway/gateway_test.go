package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckSourceIsAllowed(t *testing.T) {
	var gw = gatewayFrontend{
		workers:        nil,
		allowedSources: []string{}, // no restrictions
	}

	err := gw.checkSourceIsAllowed("anything")
	assert.NoError(t, err)

	gw.allowedSources = []string{"docker-registry.wikimedia.org/repos/releng/blubber/buildkit"}
	err = gw.checkSourceIsAllowed("docker-registry.wikimedia.org/repos/releng/blubber/buildkit")
	assert.NoError(t, err)
	err = gw.checkSourceIsAllowed("docker-registry.wikimedia.org/repos/releng/blubber/buildkit:v1.2.3")
	assert.NoError(t, err)
	err = gw.checkSourceIsAllowed("docker-registry.wikimedia.org/something-else")
	assert.Error(t, err)

	gw.allowedSources = []string{"implicit-docker-io-reference", "docker/dockerfile"}
	// This source will be rejected because after parsing it becomes
	// "docker.io/library/implicit-docker-io-reference" which does not match
	// the allowed source of "implicit-docker-io-reference".
	err = gw.checkSourceIsAllowed("implicit-docker-io-reference")
	assert.Error(t, err)
	// "docker/dockerfile" expands to "docker.io/docker/dockerfile", so no match here.
	err = gw.checkSourceIsAllowed("docker/dockerfile")
	assert.Error(t, err)

	gw.allowedSources = []string{"docker-registry.wikimedia.org/*"}
	err = gw.checkSourceIsAllowed("docker-registry.wikimedia.org/something-else")
	assert.NoError(t, err)
	err = gw.checkSourceIsAllowed("docker-registry.wikimedia.org/topdir/below")
	assert.Error(t, err)

	gw.allowedSources = []string{"docker-registry.wikimedia.org/**"}
	err = gw.checkSourceIsAllowed("docker-registry.wikimedia.org/topdir/below")
	assert.NoError(t, err)
}
