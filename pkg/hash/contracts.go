package hash

type (
	Hasher interface {
		Generate(ent string) (string, error)
		Compare(hashed, ent string) error
	}
)
