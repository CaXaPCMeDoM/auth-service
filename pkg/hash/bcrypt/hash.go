package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	cost int
}

func NewBcryptHasher(cost int) *Hasher {
	return &Hasher{
		cost: cost,
	}
}

func (h *Hasher) Generate(ent string) (string, error) {
	if len(ent) > 72 {
		ent = ent[:72]
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(ent), h.cost)

	return string(hash), err
}

func (h *Hasher) Compare(hashed, ent string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(ent))
}
