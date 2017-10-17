package errs

import "errors"

var (
    ErrBadFormat = errors.New("BadFromat")
    ErrNotVerifyTicket = errors.New("NotVerifyTicket")
)
