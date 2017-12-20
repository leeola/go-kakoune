package api

import (
	"fmt"
	"strconv"
)

const (
	var_prefix = "kak_"
	// NOTE(leeola): this is combined with var prefix in the method,
	// so it doesn't need to be prefixed in the string.
	opt_prefix = "opt_"
)

func (k *Kak) Option(key string) (string, error) {
	return k.Var(opt_prefix + key)
}

func (k *Kak) Var(key string) (string, error) {
	v, ok := k.funcVars[var_prefix+key]
	if !ok {
		// TODO(leeola): check the current commands to see if the given var
		// was even specified, so a more informative error can be returned to
		// the user.

		return "", fmt.Errorf("var not available: %q", key)
	}

	return v, nil
}

func (k *Kak) VarInt(key string) (int, error) {
	v, ok := k.funcVars[var_prefix+key]
	if !ok {
		// TODO(leeola): check the current commands to see if the given var
		// was even specified, so a more informative error can be returned to
		// the user.

		return 0, fmt.Errorf("var not available: %q", key)
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (k *Kak) Arg(i int) (string, error) {
	if i > len(k.funcArgs) {
		return "", fmt.Errorf("argument not given: %d", i)
	}

	return k.funcArgs[i], nil
}
