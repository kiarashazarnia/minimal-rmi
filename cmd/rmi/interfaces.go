package interfaces

type Hello interface {
	SayHello() string
}

type Calculator interface {
	Sum(a float32, b float32) float32
	Subtract(a float32, b float32) float32
	Multiple(a float32, b float32) float32
	Devide(a float32, b float32) float32
}
