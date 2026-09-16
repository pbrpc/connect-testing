//revive:disable:package-comments
package addr

// Mock implements net.Addr for testing.
type Mock struct {
	network string
	address string
}

// New creates a new mock address with the given address string.
func New(address string) *Mock {
	return &Mock{
		network: "tcp",
		address: address,
	}
}

// NewWithNetwork creates a new mock address with custom network and address.
func NewWithNetwork(network, address string) *Mock {
	return &Mock{
		network: network,
		address: address,
	}
}

// Network answers with the address's network name, "tcp" unless
// NewWithNetwork gave another.
func (m *Mock) Network() string {
	return m.network
}

// String answers with the address as given.
func (m *Mock) String() string {
	return m.address
}
