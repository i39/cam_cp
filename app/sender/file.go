package sender

import "context"

type File struct {
	Dir string
}

func (f *File) Run(ctx context.Context) error {
	return nil
}
