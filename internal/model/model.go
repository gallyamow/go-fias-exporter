package model

type FileInfo struct {
	Path string
	Size int64
	Hash string
	Err  error
}

func (f FileInfo) WithErr(err error) FileInfo {
	f.Err = err
	return f
}

func (f FileInfo) WithHash(hash string) FileInfo {
	f.Hash = hash
	return f
}
