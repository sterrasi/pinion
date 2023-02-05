package db

type IsoLevel string

// Transaction isolation levels
const (
	Serializable    IsoLevel = "serializable"
	RepeatableRead  IsoLevel = "repeatable read"
	ReadCommitted   IsoLevel = "read committed"
	ReadUncommitted IsoLevel = "read uncommitted"
)

// AccessMode is the transaction access mode (read write or read only)
type AccessMode string

// Transaction access modes
const (
	ReadWrite AccessMode = "read write"
	ReadOnly  AccessMode = "read only"
)

// DeferrableMode is the transaction deferrable mode (deferrable or not deferrable)
type DeferrableMode string

// Transaction deferrable modes
const (
	Deferrable    DeferrableMode = "deferrable"
	NotDeferrable DeferrableMode = "not deferrable"
)

// TransactionOptions are transaction modes within a transaction block
type TransactionOptions struct {
	IsoLevel       IsoLevel
	AccessMode     AccessMode
	DeferrableMode DeferrableMode
}
