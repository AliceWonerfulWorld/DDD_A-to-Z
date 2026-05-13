package security

type IDGenerator struct {
	prefix string
	tokens *SecureTokenGenerator
}

func NewIDGenerator(prefix string) *IDGenerator {
	return &IDGenerator{
		prefix: prefix,
		tokens: NewSecureTokenGenerator(),
	}
}

func (g *IDGenerator) NewID() (string, error) {
	token, err := g.tokens.NewToken()
	if err != nil {
		return "", err
	}

	return g.prefix + "_" + token, nil
}
