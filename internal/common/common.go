package common

import "time"

var (
	ConsumerReadTimeout   = time.Duration(-1)
	ConsumerHandleTimeout = time.Duration(5 * time.Second)
	CreateLedgerTopic     = "create_ledger"
)
