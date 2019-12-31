package face

type Issh interface {
	Run(cmd string) (string, error)
}
