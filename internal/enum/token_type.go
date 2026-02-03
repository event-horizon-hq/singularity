package enum

type TokenType string

const (
	TokenTypeMaster TokenType = "master"
	TokenTypeSlave  TokenType = "slave"
)

func (t TokenType) IsValid() bool {
	switch t {
	case TokenTypeMaster, TokenTypeSlave:
		return true
	default:
		return false
	}
}
