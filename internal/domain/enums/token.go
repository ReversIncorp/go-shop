package enums

type Token int

const (
	Access  Token = iota // 0
	Refresh              // 1
)

func (t Token) String() string {
	return [...]string{"access_token", "refresh_token"}[t]
}
