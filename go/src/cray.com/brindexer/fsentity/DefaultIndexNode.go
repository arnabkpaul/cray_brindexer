package fsentity

import (
	"bufio"
	"cray.com/brindexer/utils"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func IndexDirName() string {
	return dbName
}

type DefaultIndexNode struct {
	baseDir  string
	indexDir *string
}

func NewDefaultIndexNode(dir string, indexDir *string) DefaultIndexNode {
	dir = strings.TrimSuffix(dir, string(os.PathSeparator))
	conn := DefaultIndexNode{dir, indexDir}
	conn.init()
	return conn
}

func (node *DefaultIndexNode) init() {

	os.MkdirAll(node.DBDir(), 0777)
	//lfs setstripe --pool flash /path/to/db
	args := []string{"setstripe", "-pool", "flash", node.DBDir()}
	cmd := exec.Command("lfs", args...)
	cmd.Start()
	cmd.Wait()
}

func (node *DefaultIndexNode) BaseDir() string {
	return node.baseDir
}

func (node *DefaultIndexNode) RPath(path string) string {

	rpath, err := filepath.Rel(node.baseDir, path)
	if err != nil {
		return ""
	}
	return rpath
}

func (node *DefaultIndexNode) makeCustomeDbDir() string {
	//If we don't store index dbs in-tree, we use the md5 as part of the
	//db path to avoid conflicts
	pmd5 := utils.Md5FromString(node.baseDir)
	pname := filepath.Base(node.baseDir)
	return path.Join(*node.indexDir, pname+"_"+pmd5, dbName)
}

func (node *DefaultIndexNode) DBDir() string {
	var dbDir string
	if node.indexDir == nil {
		dbDir = path.Join(node.baseDir, dbName)
	} else {
		dbDir = node.makeCustomeDbDir()
	}
	//fmt.Println(dbDir)
	return dbDir
}

func (node *DefaultIndexNode) DBBase() *string {
	return node.indexDir
}

func (node *DefaultIndexNode) DBCount() int {
	//Should read from index info
	return TotalDbCount
}

func (node *DefaultIndexNode) WriterCount() int {
	//Should read from index info
	return DbWriterCount
}

func (node *DefaultIndexNode) DBCountPerWriter() int {
	//Should read from index info
	return node.DBCount() / node.WriterCount()
}

func (node *DefaultIndexNode) Number2FileName(num int) string {
	return fmt.Sprintf("%02x", num%node.DBCountPerWriter())
}

func (node *DefaultIndexNode) Number2WriterName(num int) string {
	//This is not really dir name, it's first 2 letter of the db file name
	//We keep this for flexibility in case we want to use dir later
	return fmt.Sprintf("%02x", num%node.WriterCount())
}

func (node *DefaultIndexNode) WriterKeyFromPathMd5(md5 string) int {
	//First 2 digits as dir name
	str := md5[0:2]
	key, _ := strconv.ParseInt(str, 16, 32)
	//fmt.Println("md5 %s, Writter key %d", md5, key)
	return int(key) % node.WriterCount()
}

func (node *DefaultIndexNode) DBWrapperKeyFromPathMd5(md5 string) int {
	//next 2 letter digits as file name
	str := md5[2:4]
	key, _ := strconv.ParseInt(str, 16, 32)
	//fmt.Println("md5 %s, wrapper key %d", md5, key)
	return int(key) % node.DBCountPerWriter()
}

func (node *DefaultIndexNode) DBFileFromPathMd5(md5 string) string {
	writerKey := node.WriterKeyFromPathMd5(md5)
	dbKeyKey := node.DBWrapperKeyFromPathMd5(md5)
	writerName := node.Number2WriterName(writerKey)
	dbFileName := writerName + node.Number2FileName(dbKeyKey) + ".db"
	dbFile := filepath.Join(node.DBDir(), dbFileName)
	return dbFile
}

func (node *DefaultIndexNode) PathInfo(path string) *PathInfo {
	rpath, err := filepath.Rel(node.baseDir, path)
	if err != nil {
		return nil
	}
	md5 := utils.Md5FromString(rpath)
	return &PathInfo{rpath: rpath, pathmd5: md5}
}

func (node *DefaultIndexNode) IsSubIndexNode(path string) bool {
	if path == node.baseDir {
		return false
	}
	dbdir := filepath.Join(path, IndexDirName())
	_, err := os.Stat(dbdir)
	return !os.IsNotExist(err)
}

func writeNumberToFile(outfile string, number int) {
	f, _ := os.Create(outfile)
	defer f.Close()
	//-1 is to be compatible with python tools.
	f.WriteString(strconv.Itoa(number))
}

func (node *DefaultIndexNode) SaveLevel(level int) {
	lfile := path.Join(node.DBDir(), levelFile)
	writeNumberToFile(lfile, level)
}

func readNumberFromFile(infile string) int {
	file, err := os.Open(infile)
	defer file.Close()
	if err != nil {
		return -1
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	lstr := scanner.Text()
	if err := scanner.Err(); err != nil {
		return -1
	}
	number, err := strconv.Atoi(lstr)
	if err != nil {
		return -1
	}
	//fmt.Println("Read level:", level)
	return number
}
func (node *DefaultIndexNode) ReadLevel() int {
	lfile := path.Join(node.DBDir(), levelFile)
	return readNumberFromFile(lfile)
}

func (node *DefaultIndexNode) SaveLastScan() {

	sfile := path.Join(node.DBDir(), lastScanFile)

	writeNumberToFile(sfile, time.Now().Second())
}
func (node *DefaultIndexNode) ReadLastScan() int {
	sfile := path.Join(node.DBDir(), lastScanFile)
	return readNumberFromFile(sfile)
}

func (node *DefaultIndexNode) Close() {
	//PLace holder, if the index is no in-tree, copy the index to base
}
