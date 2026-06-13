package urlrepo

import (
	"context"

	"github.com/sanctumlabs/curtz/app/internal/domain/url"
)

// ---------------------------------------------------------------------------
// Helper: injectMockUrlWriteTx
// Wires a mockUrlWriteQuerier into the adapter, replacing the real DB
// transaction executor. The fn passed to withTx is called directly with
// the mock querier — no real connection or transaction is involved.
// ---------------------------------------------------------------------------

func injectMockUrlWriteTx(adapter *urlWriteRepositoryAdapter, q UrlWriteQuerier) {
	adapter.withTx = func(ctx context.Context, fn func(UrlWriteQuerier) (url.URL, error)) (url.URL, error) {
		return fn(q)
	}
}
