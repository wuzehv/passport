package file

import "os"

func AppendOpenFile(logDir, filename string) (*os.File, error) {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err = os.MkdirAll(logDir, 0777); err != nil {
			return nil, err
		}
	}

	finalFile := logDir + string(os.PathSeparator) + filename

	f, err := os.OpenFile(finalFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return f, nil
}
