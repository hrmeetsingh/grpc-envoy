package port

import "time"

// Clock abstracts time retrieval (outbound port).
type Clock interface {
	Now() time.Time
}
