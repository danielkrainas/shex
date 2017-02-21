package fsutils

func copyFile(src string, dst string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return -1, err
	}

	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return -1, err
	}

	read, err := io.Copy(dstFile, srcFile)
	if err != nil {
		dstFile.Close()
		return -1, err
	}

	err = dstFile.Close()
	return read, err
}

func dirExists(filePath string) bool {
	stat, err := os.Stat(filePath)
	if err != nil || !stat.IsDir() {
		return false
	}

	return true
}

func clearDirectory(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		p := filepath.Join(dirPath, f.Name())
		if f.IsDir() {
			err = os.RemoveAll(p)
		} else {
			err = os.Remove(p)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return true
}
