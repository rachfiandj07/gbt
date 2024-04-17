package state

type (
	Wrapper struct {
		LocalUrl string
	}
)

var State Wrapper

func init() {

}

func InitWithConfig() {
	State.LocalUrl = ""
}
