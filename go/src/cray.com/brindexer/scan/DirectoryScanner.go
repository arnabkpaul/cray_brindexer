package scan

import (
	_ "fmt"
	"github.com/karrick/godirwalk"
	"os"
	"path/filepath"
	"strings"
)

type ScanCallBack interface {
	ProcessRootDir(path string, name string) error
	ProcessLeafDir(path string, name string) error
	ProcessMidDir(path string, name string) error
	ProcessMidFile(path string) error
}

type DirectoryScanner struct {
	//indexNode fsentity.IndexNode
	//Lowest level to scan to
	//0, (noLeaf), all dirs under root are middle dir, full scan
	level int
}

func NewDirectoryScanner(level int) *DirectoryScanner {

	return &DirectoryScanner{level: level}
}

func (s *DirectoryScanner) Scan(target string, callback ScanCallBack) int64 {
	//fmt.Println("Scan level ", s.level)
	target = strings.TrimSuffix(target, string(os.PathSeparator))
	count := int64(0)
	godirwalk.Walk(target, &godirwalk.Options{
		Unsorted:            true,
		FollowSymbolicLinks: false,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			count++
			//fmt.Println("Entry:", osPathname)
			if de.IsDir() {
				rpath, err := filepath.Rel(target, osPathname)
				if err != nil {
					return filepath.SkipDir
				}
				//fmt.Println("Rel Path:", rpath)
				if rpath == "." {
					return callback.ProcessRootDir(osPathname, de.Name())
				}
				//All middle
				if s.level == 0 {
					return callback.ProcessMidDir(osPathname, de.Name())
				}
				//Note: a dir under root has zero slash in rpath, it is level 1
				slashes := strings.Count(rpath, string(os.PathSeparator))
				//fmt.Println("Rpath, slashes:", rpath, slashes)
				if (slashes + 1) == s.level {
					return callback.ProcessLeafDir(osPathname, de.Name())
				}
				return callback.ProcessMidDir(osPathname, de.Name())

			} else if !de.IsDevice() {
				callback.ProcessMidFile(osPathname)
			}

			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
	})
	//fmt.Println("Total file system objects", count)

	return count

}
