package file

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gabriel-vasile/mimetype"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gorm.io/gorm"
)

var (
	ErrFileIsNotValid     = errors.New("file is not valid")
	ErrFileNameIsNotValid = errors.New("file name is not valid")
	ErrHashInit           = errors.New("new hash failed")
	ErrFileIsNotOpen      = errors.New("file is not open")
	ErrFileAlreadyOpen    = errors.New("file already open")
	ErrNoContentTypeFound = errors.New("no content type found")
	ErrNoContentTypeMatch = errors.New("no content type match")
)

// File represents the metadata of a file
type File struct {
	database.Model
	Path        string         `json:"path,omitempty"`
	URL         *database.URL  `json:"url,omitempty"`
	Title       string         `json:"title,omitempty"`
	Description string         `json:"description,omitempty"`
	Tags        string         `json:"tags,omitempty"`
	Name        string         `json:"name,omitempty"`
	Size        int64          `json:"size,omitempty"`
	ModTime     time.Time      `json:"mod_time,omitempty"`
	ContentType string         `json:"content_type,omitempty"`
	Hash        string         `json:"hash,omitempty" gorm:"uniqueIndex"`
	SyncedAt    *database.Time `json:"synced_at,omitempty" gorm:"index"`
	OwnedBy     any            `json:"owned_by,omitempty" gorm:"-:all"`
	Stored      *File          `json:"stored,omitempty" gorm:"-:all"`
	// private
	err error
	src *os.File
}

func NewFile(path string, reContentType *regexp.Regexp, reName *regexp.Regexp) *File {
	f := &File{}

	if info, err := os.Lstat(path); err != nil {
		f.err = err
		return f
	} else {
		f.Name = info.Name()
		f.Size = info.Size()
		f.ModTime = info.ModTime()
		f.Path = path
	}

	if err := f.Open().
		MatchFileName(reName).
		MatchContentTypeFile(reContentType).
		GenerateHash().
		Close().err; err != nil {
		f.err = err
		return f
	} else {
		return f
	}
}

func HashValid(hash string) bool {
	hash = strings.TrimSpace(hash)
	return hash != ""
}

func (f *File) BeforeCreate(tx *gorm.DB) error {
	t := database.NowTime()
	f.SyncedAt = &t
	return nil
}

func (f *File) Synced() bool {
	return f.SyncedAt != nil
}

func (f *File) HashValid() bool {
	return HashValid(f.Hash)
}

func (f *File) String() string {
	return fmt.Sprintf("File: name(%s) path (%s) size (%s)", f.Name, f.Path, humanize.Bytes(uint64(f.Size)))
}

func (f *File) Error() error {
	return f.err
}

func (f *File) IsOpen() error {
	if f.src != nil {
		return nil
	} else {
		return ErrFileIsNotOpen
	}
}

func (f *File) Open() *File {
	if f.src != nil || f.err != nil {
		return f
	} else if file, err := os.Open(f.Path); err != nil {
		f.err = err
		return f
	} else {
		f.src = file
		return f
	}
}

func (f *File) Close() *File {
	if f.err != nil {
		return f
	} else {
		f.src.Close()
		return f
	}
}

func (f *File) IsValid() error {
	// TODO replace len(f.Name) > 0 with regex
	if f.Hash == "" ||
		f.Path == "" ||
		f.Size == 0 ||
		f.Name == "" {
		return ErrFileIsNotValid
	} else {
		return nil
	}
}

// GenerateHash generates a unique ID for the file based on its path and modification time.
func (f *File) GenerateHash() *File {
	if f.err != nil {
		return f
	} else if err := f.IsOpen(); err != nil {
		f.err = ErrFileIsNotOpen
		return f
	} else if hash := md5.New(); hash == nil {
		f.err = ErrHashInit
		return f
	} else if _, err := io.Copy(hash, f.src); err != nil {
		f.err = err
		return f
	} else {
		f.Hash = hex.EncodeToString(hash.Sum(nil))
		return f
	}
}

// MatchContentTypeFile ...
func (f *File) MatchContentTypeFile(reContentType *regexp.Regexp) *File {
	if reContentType == nil {
		return f
	} else if err := f.IsOpen(); err != nil {
		f.err = err
		return f
	} else if mt, err := mimetype.DetectReader(f.src); err != nil {
		f.err = err
		return f
	} else if f.ContentType = mt.String(); f.ContentType == "" {
		f.err = ErrNoContentTypeFound
		return f
	} else if !reContentType.MatchString(mt.String()) {
		f.err = ErrNoContentTypeMatch
		return f
	} else {
		return f
	}
}

func (f *File) MatchFileName(re *regexp.Regexp) *File {
	if re == nil {
		return f
	} else if !re.MatchString(f.Name) {
		f.err = ErrFileNameIsNotValid
		return f
	} else {
		return f
	}
}
