// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

var initFlagNames map[int64]string
var releaseFlagNames map[int64]string
var OpenFlagNames map[int64]string
var FuseOpenFlagNames map[int64]string
var accessFlagName map[int64]string
var writeFlagNames map[int64]string
var readFlagNames map[int64]string

func init() {
	writeFlagNames = map[int64]string{
		WRITE_CACHE:     "CACHE",
		WRITE_LOCKOWNER: "LOCKOWNER",
	}
	readFlagNames = map[int64]string{
		READ_LOCKOWNER: "LOCKOWNER",
	}
	initFlagNames = map[int64]string{
		CAP_ASYNC_READ:         "ASYNC_READ",
		CAP_POSIX_LOCKS:        "POSIX_LOCKS",
		CAP_FILE_OPS:           "FILE_OPS",
		CAP_ATOMIC_O_TRUNC:     "ATOMIC_O_TRUNC",
		CAP_EXPORT_SUPPORT:     "EXPORT_SUPPORT",
		CAP_BIG_WRITES:         "BIG_WRITES",
		CAP_DONT_MASK:          "DONT_MASK",
		CAP_SPLICE_WRITE:       "SPLICE_WRITE",
		CAP_SPLICE_MOVE:        "SPLICE_MOVE",
		CAP_SPLICE_READ:        "SPLICE_READ",
		CAP_FLOCK_LOCKS:        "FLOCK_LOCKS",
		CAP_IOCTL_DIR:          "IOCTL_DIR",
		CAP_AUTO_INVAL_DATA:    "AUTO_INVAL_DATA",
		CAP_READDIRPLUS:        "READDIRPLUS",
		CAP_READDIRPLUS_AUTO:   "READDIRPLUS_AUTO",
		CAP_ASYNC_DIO:          "ASYNC_DIO",
		CAP_WRITEBACK_CACHE:    "WRITEBACK_CACHE",
		CAP_NO_OPEN_SUPPORT:    "NO_OPEN_SUPPORT",
		CAP_PARALLEL_DIROPS:    "PARALLEL_DIROPS",
		CAP_POSIX_ACL:          "POSIX_ACL",
		CAP_HANDLE_KILLPRIV:    "HANDLE_KILLPRIV",
		CAP_ABORT_ERROR:        "ABORT_ERROR",
		CAP_MAX_PAGES:          "MAX_PAGES",
		CAP_CACHE_SYMLINKS:     "CACHE_SYMLINKS",
		CAP_NO_OPENDIR_SUPPORT: "NO_OPENDIR_SUPPORT",
	}
	releaseFlagNames = map[int64]string{
		RELEASE_FLUSH: "FLUSH",
	}
	OpenFlagNames = map[int64]string{
		int64(os.O_WRONLY):        "WRONLY",
		int64(os.O_RDWR):          "RDWR",
		int64(os.O_APPEND):        "APPEND",
		int64(syscall.O_ASYNC):    "ASYNC",
		int64(os.O_CREATE):        "CREAT",
		int64(os.O_EXCL):          "EXCL",
		int64(syscall.O_NOCTTY):   "NOCTTY",
		int64(syscall.O_NONBLOCK): "NONBLOCK",
		int64(os.O_SYNC):          "SYNC",
		int64(os.O_TRUNC):         "TRUNC",

		int64(syscall.O_CLOEXEC):   "CLOEXEC",
		int64(syscall.O_DIRECTORY): "DIRECTORY",
	}
	FuseOpenFlagNames = map[int64]string{
		FOPEN_DIRECT_IO:   "DIRECT",
		FOPEN_KEEP_CACHE:  "CACHE",
		FOPEN_NONSEEKABLE: "NONSEEK",
	}
	accessFlagName = map[int64]string{
		X_OK: "x",
		W_OK: "w",
		R_OK: "r",
	}
}

func FlagString(names map[int64]string, fl int64, def string) string {
	s := []string{}
	for k, v := range names {
		if fl&k != 0 {
			s = append(s, v)
			fl ^= k
		}
	}
	if len(s) == 0 && def != "" {
		s = []string{def}
	}
	if fl != 0 {
		s = append(s, fmt.Sprintf("0x%x", fl))
	}

	return strings.Join(s, ",")
}

func (me *ForgetIn) string() string {
	return fmt.Sprintf("{Nlookup=%d}", me.Nlookup)
}

func (me *_BatchForgetIn) string() string {
	return fmt.Sprintf("{Count=%d}", me.Count)
}

func (me *MkdirIn) string() string {
	return fmt.Sprintf("{0%o (0%o)}", me.Mode, me.Umask)
}

func (me *Rename1In) string() string {
	return fmt.Sprintf("{i%d}", me.Newdir)
}

func (me *RenameIn) string() string {
	return fmt.Sprintf("{i%d %x}", me.Newdir, me.Flags)
}

func (me *SetAttrIn) string() string {
	s := []string{}
	if me.Valid&FATTR_MODE != 0 {
		s = append(s, fmt.Sprintf("mode 0%o", me.Mode))
	}
	if me.Valid&FATTR_UID != 0 {
		s = append(s, fmt.Sprintf("uid %d", me.Uid))
	}
	if me.Valid&FATTR_GID != 0 {
		s = append(s, fmt.Sprintf("gid %d", me.Gid))
	}
	if me.Valid&FATTR_SIZE != 0 {
		s = append(s, fmt.Sprintf("size %d", me.Size))
	}
	if me.Valid&FATTR_ATIME != 0 {
		s = append(s, fmt.Sprintf("atime %d.%09d", me.Atime, me.Atimensec))
	}
	if me.Valid&FATTR_MTIME != 0 {
		s = append(s, fmt.Sprintf("mtime %d.%09d", me.Mtime, me.Mtimensec))
	}
	if me.Valid&FATTR_FH != 0 {
		s = append(s, fmt.Sprintf("fh %d", me.Fh))
	}
	// TODO - FATTR_ATIME_NOW = (1 << 7), FATTR_MTIME_NOW = (1 << 8), FATTR_LOCKOWNER = (1 << 9)
	return fmt.Sprintf("{%s}", strings.Join(s, ", "))
}

func (me *ReleaseIn) string() string {
	return fmt.Sprintf("{Fh %d %s %s L%d}",
		me.Fh, FlagString(OpenFlagNames, int64(me.Flags), ""),
		FlagString(releaseFlagNames, int64(me.ReleaseFlags), ""),
		me.LockOwner)
}

func (me *OpenIn) string() string {
	return fmt.Sprintf("{%s}", FlagString(OpenFlagNames, int64(me.Flags), "O_RDONLY"))
}

func (me *OpenOut) string() string {
	return fmt.Sprintf("{Fh %d %s}", me.Fh,
		FlagString(FuseOpenFlagNames, int64(me.OpenFlags), ""))
}

func (me *InitIn) string() string {
	return fmt.Sprintf("{%d.%d Ra 0x%x %s}",
		me.Major, me.Minor, me.MaxReadAhead,
		FlagString(initFlagNames, int64(me.Flags), ""))
}

func (me *InitOut) string() string {
	return fmt.Sprintf("{%d.%d Ra 0x%x %s %d/%d Wr 0x%x Tg 0x%x}",
		me.Major, me.Minor, me.MaxReadAhead,
		FlagString(initFlagNames, int64(me.Flags), ""),
		me.CongestionThreshold, me.MaxBackground, me.MaxWrite,
		me.TimeGran)
}

func (s *FsyncIn) string() string {
	return fmt.Sprintf("{Fh %d Flags %x}", s.Fh, s.FsyncFlags)
}

func (me *SetXAttrIn) string() string {
	return fmt.Sprintf("{sz %d f%o}", me.Size, me.Flags)
}

func (me *GetXAttrIn) string() string {
	return fmt.Sprintf("{sz %d}", me.Size)
}

func (me *GetXAttrOut) string() string {
	return fmt.Sprintf("{sz %d}", me.Size)
}

func (me *AccessIn) string() string {
	return fmt.Sprintf("{u=%d g=%d %s}",
		me.Uid,
		me.Gid,
		FlagString(accessFlagName, int64(me.Mask), ""))
}

func (me *FlushIn) string() string {
	return fmt.Sprintf("{Fh %d}", me.Fh)
}

func (me *AttrOut) string() string {
	return fmt.Sprintf(
		"{tA=%gs %v}",
		ft(me.AttrValid, me.AttrValidNsec), &me.Attr)
}

// ft converts (seconds , nanoseconds) -> float(seconds)
func ft(tsec uint64, tnsec uint32) float64 {
	return float64(tsec) + float64(tnsec)*1E-9
}

// Returned by LOOKUP
func (me *EntryOut) string() string {
	return fmt.Sprintf("{i%d g%d tE=%gs tA=%gs %v}",
		me.NodeId, me.Generation, ft(me.EntryValid, me.EntryValidNsec),
		ft(me.AttrValid, me.AttrValidNsec), &me.Attr)
}

func (me *CreateOut) string() string {
	return fmt.Sprintf("{i%d g%d %v %v}", me.NodeId, me.Generation, &me.EntryOut, &me.OpenOut)
}

func (me *StatfsOut) string() string {
	return fmt.Sprintf(
		"{blocks (%d,%d)/%d files %d/%d bs%d nl%d frs%d}",
		me.Bfree, me.Bavail, me.Blocks, me.Ffree, me.Files,
		me.Bsize, me.NameLen, me.Frsize)
}

func (o *NotifyInvalEntryOut) string() string {
	return fmt.Sprintf("{parent i%d sz %d}", o.Parent, o.NameLen)
}

func (o *NotifyInvalInodeOut) string() string {
	return fmt.Sprintf("{i%d [%d +%d)}", o.Ino, o.Off, o.Length)
}

func (o *NotifyInvalDeleteOut) string() string {
	return fmt.Sprintf("{parent i%d ch i%d sz %d}", o.Parent, o.Child, o.NameLen)
}

func (o *NotifyStoreOut) string() string {
	return fmt.Sprintf("{i%d [%d +%d)}", o.Nodeid, o.Offset, o.Size)
}

func (o *NotifyRetrieveOut) string() string {
	return fmt.Sprintf("{> %d: i%d [%d +%d)}", o.NotifyUnique, o.Nodeid, o.Offset, o.Size)
}

func (i *NotifyRetrieveIn) string() string {
	return fmt.Sprintf("{[%d +%d)}", i.Offset, i.Size)
}

func (f *FallocateIn) string() string {
	return fmt.Sprintf("{Fh %d [%d +%d) mod 0%o}",
		f.Fh, f.Offset, f.Length, f.Mode)
}

func (f *LinkIn) string() string {
	return fmt.Sprintf("{Oldnodeid: %d}", f.Oldnodeid)
}

func (o *WriteOut) string() string {
	return fmt.Sprintf("{%db }", o.Size)

}
func (i *CopyFileRangeIn) string() string {
	return fmt.Sprintf("{Fh %d [%d +%d) => i%d Fh %d [%d, %d)}",
		i.FhIn, i.OffIn, i.Len, i.NodeIdOut, i.FhOut, i.OffOut, i.Len)
}

func (in *InterruptIn) string() string {
	return fmt.Sprintf("{ix %d}", in.Unique)
}

var seekNames = map[uint32]string{
	0: "SET",
	1: "CUR",
	2: "END",
	3: "DATA",
	4: "HOLE",
}

func (in *LseekIn) string() string {
	return fmt.Sprintf("{Fh %d [%s +%d)}", in.Fh,
		seekNames[in.Whence], in.Offset)
}

func (o *LseekOut) string() string {
	return fmt.Sprintf("{%d}", o.Offset)
}

// Print pretty prints FUSE data types for kernel communication
func Print(obj interface{}) string {
	t, ok := obj.(interface {
		string() string
	})
	if ok {
		return t.string()
	}
	return fmt.Sprintf("%T: %v", obj, obj)
}
