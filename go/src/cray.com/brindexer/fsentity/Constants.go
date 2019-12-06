package fsentity

import ()

const MaxCacheSize int = 800
const TotalDbCount int = 64
const DbWriterCount = 64

const dbName string = "._dbindex_"
const levelFile string = "index_level"
const lastScanFile string = "last_scan"

const InsertFileSql string = "replace into entries_0 (" +
	"name, type, pathmd5, " +
	"inode, mode, nlink, uid, gid, size, blksize, blocks, atime, mtime, ctime, " +
	"poolname, ostindices, " +
	"mirrorstate, hsmarchid, hsmcompat, hsmflags, hsmarchver, seq)" +
	"values ('%s','%s', '%s'," +
	"'%d','%d','%d','%d','%d','%d','%d','%d','%d','%d','%d'," +
	"'%s','%s'," +
	"'%d','%d','%d','%d','%d','%d')"

const InsertPathSql string = "replace into path (path, pmd5, lastmtime, fcnt, total_size, pseq)" +
	"values ('%s', '%s', '%d', '%d', '%d', '%d')"

const FileSql string = "CREATE TABLE if not exists entries_0( " +
	"name TEXT not null, type TEXT not null, pathmd5 TEXT not null," +
	"inode INT64, mode INT64, nlink INT64, uid INT64, gid INT64, size INT64," +
	"blksize INT64, blocks INT64, atime INT64, mtime INT64, ctime INT64, " +
	"poolname TEXT, ostindices TEXT, mirrorstate INT64, hsmarchid INT64," +
	"hsmcompat INT64, hsmflags INT64, hsmarchver, seq INT64," +
	"primary key(pathmd5, name))"

const PathSql string = "create table if not exists path (" +
	"pmd5 text primary key collate nocase," +
	"path text not null collate nocase," +
	"lastmtime integer not null default 0," +
	"fcnt integer not null default 0," +
	"total_size integer not null default 0," +
	"pseq INT64 not null default 0)"

const insertPathSql string = "replace into path (path, pmd5, lastmtime, fcnt, total_size, pseq)" +
	"values ('%s', '%s', '%d', '%d', '%d', '%d')"

const deletePathSql string = "delete from path where pmd5='%s'; delete from entries_0 where pathmd5='%s'"
