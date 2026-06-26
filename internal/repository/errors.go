package repository

import "errors"

// ErrPartNotFound indica que a peça solicitada não existe no repositório.
var ErrPartNotFound = errors.New("peça não encontrada")
