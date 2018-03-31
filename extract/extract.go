package extract

import (
	"log"
	"archive/zip"
	"path/filepath"
	"strings"
	"io"
	"os"
	"github.com/samkreter/golang-utils/downloader"
)

	// fileURL, ok := os.LookupEnv("FILE_URL")
	// if (!ok){
	// 	log.Fatal("Must have FILE_URL set.")
	// }

	// filePath, ok := os.LookupEnv("FILE_PATH")
	// if (!ok){
	// 	log.Fatal("Must have FILE_URL set.")
	// }

	// unzipedDest, ok := os.LookupEnv("UNZIPED_DEST")
	// if (!ok){
	// 	unzipedDest = strings.TrimSuffix(filePath, filepath.Ext(filePath))
	// }
// fileURL := "https://ed-public-download.app.cloud.gov/downloads/CollegeScorecard_Raw_Data.zip"
// 	filePath := "./CollegeScorecard_Raw_Data.zip"

// DownloadAndUnzip takes in 
func DownloadAndUnzip(fileURL string, filePath string, unzipedDest string){
    if (unzipedDest == "") {
        unzipedDest = strings.TrimSuffix(filePath, filepath.Ext(filePath))
    }
	
	err := downloader.DownloadFile(filePath, fileURL)
	if err != nil {
		log.Fatal(err)
	}

	_, err = Unzip(filePath, unzipedDest)
	if err != nil {
		log.Fatal(err)
	}
}

// Unzip will un-compress a zip archive,
// moving all files and folders to an output directory
func Unzip(src, dest string) ([]string, error) {

    var filenames []string

    r, err := zip.OpenReader(src)
    if err != nil {
        return filenames, err
    }
    defer r.Close()

    for _, f := range r.File {

        rc, err := f.Open()
        if err != nil {
            return filenames, err
        }
        defer rc.Close()

        // Store filename/path for returning and using later on
        fpath := filepath.Join(dest, f.Name)
        filenames = append(filenames, fpath)

        if f.FileInfo().IsDir() {

            // Make Folder
            os.MkdirAll(fpath, os.ModePerm)

        } else {

            // Make File
            var fdir string
            if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
                fdir = fpath[:lastIndex]
            }

            err = os.MkdirAll(fdir, os.ModePerm)
            if err != nil {
                log.Fatal(err)
                return filenames, err
            }
            f, err := os.OpenFile(
                fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return filenames, err
            }
            defer f.Close()

            _, err = io.Copy(f, rc)
            if err != nil {
                return filenames, err
            }

        }
    }
    return filenames, nil
}