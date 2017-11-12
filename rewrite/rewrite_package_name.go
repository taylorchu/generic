package rewrite

// rewritePackageName sets current package name.
func (s *Spec) rewritePackageName(pkg *Package) error {
	for _, node := range pkg.Files {
		node.Name.Name = s.Package
	}
	return nil
}
