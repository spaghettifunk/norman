package commander

type Commander struct {
	Name    string
	Address string
	Port    int
	// configuration here...
}

func New()