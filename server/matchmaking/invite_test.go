package matchmaking

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestGenerateLink(t *testing.T) {
	link, err := GenerateLink(1)
	assert.NilError(t, err)
	assert.Assert(t, link != "")

	suuid, ok := users[1]
	assert.Assert(t, ok)
	assert.Equal(t, suuid, link)

	id, ok := links[link]
	assert.Assert(t, ok)
	assert.Equal(t, id, int64(1))
}

func TestJoinLink(t *testing.T) {
	initMock()
	link, err := GenerateLink(1)
	assert.NilError(t, err)

	err = JoinLink(link, 2)
	assert.NilError(t, err)

	_, ok := users[1]
	assert.Assert(t, !ok)

	_, ok = links[link]
	assert.Assert(t, !ok)
}
