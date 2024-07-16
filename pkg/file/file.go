package file

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"regexp"
	"time"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/gabriel-vasile/mimetype"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"
)

var (
	ErrFileIsNotValid     = errors.New("file is not valid")
	ErrHashInit           = errors.New("new hash failed")
	ErrFileIsNotOpen      = errors.New("file is not open")
	ErrFileAlreadyOpen    = errors.New("file already open")
	ErrNoContentTypeFound = errors.New("no content type found")
	ErrNoContentTypeMatch = errors.New("no content type match")
)

// File represents the metadata of a file
type File struct {
	simutils.Model
	Name        string        `json:"name,omitempty"`
	Size        int64         `json:"size,omitempty"`
	ModTime     time.Time     `json:"mod_time,omitempty"`
	ContentType string        `json:"content_type,omitempty"`
	Hash        string        `json:"hash,omitempty" gorm:"uniqueIndex"`
	Path        string        `json:"path,omitempty"`
	Stored      *File         `json:"stored,omitempty" gorm:"-:all"`
	ImageID     database.PID  `json:"image_id,omitempty" gorm:"default:null"`
	Image       *models.Image `json:"image,omitempty"`
	// private
	err error
	src *os.File
}

func NewFile(reContentType *regexp.Regexp, path string) *File {
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
		MatchContentTypeFile(reContentType).
		GenerateHash().
		Close().err; err != nil {
		f.err = err
		return f
	} else {
		return f
	}
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
	if err := f.IsOpen(); err != nil {
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

func MigrateFile(db *simutils.DBConnection) error {
	return db.DB.AutoMigrate(&models.Image{}, &File{})
}
