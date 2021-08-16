package file

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
)

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func UnZIP(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func RemoveFile(pathFile string) error {
	if CheckFileExist(pathFile) {
		return os.Remove(pathFile)
	}
	return nil
}

func CheckFileExist(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CheckFolderExist(pathRoot string) bool {
	info, err := os.Stat(pathRoot)
	if err != nil {
		fmt.Println("ERROR START")
		return false
	}
	if os.IsNotExist(err) {
		fmt.Println("Not exist")
		return false
	}
	return info.IsDir()
}

func CreateFolder(dir string) error {
	if len(dir) == 0 {
		return nil
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		return err
	} else if err != nil {
		return err
	}
	return nil
}

func RemoveSpecial(value string) string {
	if strings.Contains(value, ".") {
		var extension = path.Ext(value)
		var name = value[0 : len(value)-len(extension)]
		value = slug.MakeLang(name, "en") + extension
	} else {
		value = slug.MakeLang(value, "en")
	}
	return value
}

func ReplaceSpecial(value, addValue string) string {
	var dataStr string
	if strings.Contains(value, ".") {
		var arrVal = strings.Split(value, ".")
		var lenArr = len(arrVal)
		for i, val := range arrVal {
			if i == lenArr-1 {
				dataStr += "_" + addValue + "." + val
			} else {
				dataStr += slug.MakeLang(val, "en")
			}
		}
	} else {
		dataStr = slug.MakeLang(value, "en")
	}
	return dataStr
}

func ClearDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func RenameFile(pathOld, pathNew string) error {
	return os.Rename(pathOld, pathNew)
}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = addFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func CopyFolder(source string, dest string) (err error) {

	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			err = CopyFolder(sourcefilepointer, destinationfilepointer)
		} else {
			err = CopyFile(sourcefilepointer, destinationfilepointer)
		}
	}
	return err
}

func CopyFileToFolder(source, dest, fileName string) (err error) {
	if !CheckFolderExist(dest) {
		CreateFolder(dest)
	}

	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()
	dest += "/" + fileName
	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// func FilePathWalkDir(root string) ([]string, error) {
// 	var files []string
// 	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
// 		if !info.IsDir() {
// 			files = append(files, path)
// 		} else {
// 			var fNew, err = FilePathWalkDir(path)
// 			if err == nil {
// 				files = append(files, fNew...)
// 			}
// 		}
// 		return nil
// 	})
// 	return files, err
// }

func IOReadDirFile(root string) ([]string, error) {
	isCheck := CheckFolderExist(root)
	if !isCheck {
		fmt.Println("Not dir")
	}
	var files []string

	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		var name = file.Name()
		var fi, _ = os.Stat(filepath.Join(root, name))
		if !fi.IsDir() {
			files = append(files, name)
		}
	}
	return files, nil
}

func IOReadDirAll(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		var name = file.Name()
		files = append(files, name)
	}
	return files, nil
}

func ReadFile(fileName string) (string, error) {

	var f, err = ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	// Convert []byte to string and print to screen
	text := string(f)
	return text, nil
}

func SaveBase64ToImage(b64, pathSave string) error {
	dec, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}

	f, err := os.Create(pathSave)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	return nil
}
