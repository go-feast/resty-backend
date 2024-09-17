package closer

import (
	"github.com/rs/zerolog/log"
	"io"
	"maps"
)

type (
	Closer struct {
		forClose map[string]io.Closer
	}
)

func NewCloser() *Closer {
	return &Closer{forClose: make(map[string]io.Closer)}
}

func (c *Closer) Close() {
	for name, closer := range c.forClose {
		err := closer.Close()
		if err != nil {
			log.Err(err).Msgf("failed to close %s: %s", name, err)
		}
	}

	log.Info().Msg("all dependencies are closed")
}

func (c *Closer) AppendClosers(m map[string]io.Closer) {
	maps.Copy(c.forClose, m)
}

func (c *Closer) AppendCloser(closer *Closer) {
	maps.Copy(c.forClose, closer.forClose)
}
