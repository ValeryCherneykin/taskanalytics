package tests

type fakeStorageConfig struct {
	basePath string
}

func (f *fakeStorageConfig) Path() string {
	return f.basePath
}
