package identityapp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
)

// Refresh exchanges a valid refresh token for a fresh token pair.
//
// The user is re-read so that a token belonging to a since-deleted account stops working.
func (svc *Service) Refresh(ctx context.Context, refreshToken string) (TokenPair, error) {
	handlerLogPrefix := fmt.Sprintf("%s<Refresh>", svc.logPrefix)

	if refreshToken == "" {
		return TokenPair{}, errdefs.Unauthorized(errdefs.ErrTokenRequired)
	}

	userID, authErr := svc.tokens.Authenticate(refreshToken)
	if authErr != nil {
		slog.WarnContext(ctx, fmt.Sprintf("%s Refresh rejected: invalid token", handlerLogPrefix), "error", authErr)
		return TokenPair{}, errdefs.Unauthorized(errdefs.ErrTokenInvalid)
	}

	if _, fetchErr := svc.users.FetchById(ctx, userID); fetchErr != nil {
		slog.WarnContext(ctx, fmt.Sprintf("%s Refresh rejected: user no longer exists", handlerLogPrefix), "userId", userID)
		return TokenPair{}, errdefs.Unauthorized(errdefs.ErrTokenInvalid)
	}

	tokens, tokenErr := svc.issueTokens(userID)
	if tokenErr != nil {
		return TokenPair{}, tokenErr
	}

	return tokens, nil
}
