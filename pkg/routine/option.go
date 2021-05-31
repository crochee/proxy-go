package routine

type option struct {
	recoverFunc func(interface{})
}

// Recover register to pool
func Recover(f func(interface{})) func(*option) {
	return func(o *option) { o.recoverFunc = f }
}
