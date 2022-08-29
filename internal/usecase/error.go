package usecase

import "errors"

var ErrBlockNotFinalize = errors.New("block is not finalized")
var ErrConnectionProblem = errors.New("problem with connection or network")
