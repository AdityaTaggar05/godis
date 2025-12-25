package types

type RedisWrongTypeError struct{}
type RedisWrongNumArgsError struct {
	cmd string
}
type RedisInvalidStreamIDError struct{}
type RedisMinimumStreamIDError struct{}

func panicRedisWrongType() {
	panic(RedisWrongTypeError{})
}

func panicRedisWrongNumArgs(cmd string) {
	panic(RedisWrongNumArgsError{cmd: cmd})
}

func panicRedisInvalidStreamID() {
	panic(RedisInvalidStreamIDError{})
}

func panicRedisMinimumStreamID() {
	panic(RedisMinimumStreamIDError{})
}
