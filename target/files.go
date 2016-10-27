package target

type Files struct {
	Type    string // Pseudo MIME-type describing contents of tarball
	Tarball []byte // Tar'ed and gzip'ed files
}
