package chord

import (
	"crypto/sha1"
	"hash"
)

type Config struct {
	Hash     func() hash.Hash
	HashSize int
}

func DefaultConfig() *Config {
	defCng := &Config{
		Hash: sha1.New,
	}
	defCng.HashSize = defCng.Hash().Size() * 8
	return defCng
}
