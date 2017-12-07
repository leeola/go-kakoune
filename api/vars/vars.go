package vars

const (
	opt_prefix = "kak_opt_"
	BufName    = "kak_bufname"
	BufFile    = "kak_buffile"
)

func Option(key string) string {
	return opt_prefix + key
}
