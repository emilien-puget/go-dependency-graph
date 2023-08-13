package config

import "path/filepath"

const DefaultOutOfPackageDirectory = string(filepath.Separator) + "testdata" + string(filepath.Separator) + "mocks"

type Config struct {
	InPackage                  bool
	OutOfPackageMocksDirectory string
}
