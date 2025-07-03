// Agent Channel is a mutex
package channel

type Channel struct {
	mtx chan struct{}
}
