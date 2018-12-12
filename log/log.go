package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	levelfatal = 1 << iota
	levelerror
	levelwarn
	levelinfo
	leveldebug
)

const (
	LevelFatal = levelfatal
	LevelError = levelerror | LevelFatal
	LevelWarn  = levelwarn | LevelError
	LevelInfo  = levelinfo | LevelWarn
	LevelDebug = leveldebug | LevelInfo
)

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
	Lsplith                       // split by hour
	Lsplitm                       // split by minute
	Lsplitd                       // split by day
)

var logger = NewDefault()

type Logger struct {
	output     io.Writer
	level      int
	maxsize    int
	flag       int
	lastmetric int //last hour or last minute or last day
	fname      string
	sync.Mutex
	raw *log.Logger
}

func NewDefault() *Logger {
	return &Logger{
		level: LevelDebug,
		raw:   log.New(os.Stdout, "", log.Lshortfile|log.Lmicroseconds),
		flag:  Lshortfile | Lmicroseconds,
	}
}

func SetFlags(flag int) {
	logger.SetFlags(flag)
}

func SetVerbose() {
	logger.SetVerbose()
}

func SetOutPut(w io.Writer) {
	logger.SetOutPut(w)
}

func SetLevel(level int) {
	logger.SetLevel(level)
}

func SetMaxSize(max int) {
	logger.SetMaxSize(max)
}

func (l *Logger) SetFlags(flag int) {
	l.flag = flag
	l.raw.SetFlags(flag)
}

func (l *Logger) SetVerbose() {
	l.raw.SetFlags(l.flag | log.Llongfile)
}

func (l *Logger) SetOutPut(w io.Writer) {
	l.raw.SetOutput(w)
}

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) SetMaxSize(max int) {
	_, ok := l.output.(*os.File)
	if !ok {
		return
	}
	l.maxsize = max
}

func (l *Logger) SetFile(fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	l.fname = fname
	l.output = f
	l.raw.SetOutput(f)
	return nil
}

func (l *Logger) logf(format string, level int, args ...interface{}) {
	if l.level|level != l.level {
		return
	}
	prefix := "error"
	switch level {
	case leveldebug:
		prefix = "debug"
	case levelinfo:
		prefix = "info"
	case levelwarn:
		prefix = "warn"
	case levelerror:
		prefix = "error"
	case levelfatal:
		prefix = "fatal"
	}
	msg := fmt.Sprintf("[%s] %s", prefix, format)
	if len(args) != 0 {
		msg = fmt.Sprintf("%s", args...)
	}
	l.raw.Output(4, msg)
	if level == levelfatal {
		os.Exit(1)
	}
	if err := l.rotate(); err != nil {
		log.Fatal(err.Error())
	}
}

//format {fname}.{time}.{index}.{suffix}

func (l *Logger) rotate() error {
	f, ok := l.output.(*os.File)
	if !ok {
		return nil
	}
	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("stat file:%s", err.Error())
	}
	if l.fname == "" {
		return nil
	}

	if l.maxsize == 0 && l.flag < 128 {
		return nil
	}

	mtime := info.ModTime()
	name := info.Name()
	size := info.Size()
	now := time.Now()
	now_date := now.Format("2006-01-02")
	l_date := mtime.Format("2006-01-02")
	m, h := now.Minute(), now.Hour()
	lm, lh := mtime.Minute(), mtime.Hour()

	namearr := strings.Split(name, ".")
	suffix := namearr[len(namearr)-1]
	prefix := strings.Join(namearr[:len(namearr)-1], ".")

	//split by size
	if l.maxsize > 0 {
		if size > int64(l.maxsize) {
			tprefix := ""
			if l.flag&Lsplitd == Lsplitd {
				tprefix = now_date
			} else if l.flag&Lsplith == Lsplith {
				tprefix = fmt.Sprintf("%s.%d", now_date, h)
			} else if l.flag&Lsplitm == Lsplitm {
				tprefix = fmt.Sprintf("%s.%d.%d", now_date, h, m)
			}
			for {
				i := 1
				new_name := fmt.Sprintf("%s.%s.%4d.%s", prefix, tprefix, i, suffix)
				_, err := os.Open(new_name)
				if err == os.ErrNotExist {
					f.Sync()
					f.Close()
					os.Rename(name, new_name)
					f, _ = os.Create(l.fname)
					l.output = f
					l.raw.SetOutput(f)
					return nil
				}
				i++
			}
		}
	}
	//split by time
	tprefix := ""
	rotate := false
	if l.flag&Lsplitd == Lsplitd {
		if now_date != l_date {
			rotate = true
		}
		tprefix = now_date
	} else if l.flag&Lsplith == Lsplith {
		if now_date != l_date || h != lh {
			rotate = true
		}
		tprefix = fmt.Sprintf("%s.%d", now_date, h)
	} else if l.flag&Lsplitm == Lsplitm {
		if now_date != l_date || h != lh || m != lm {
			rotate = true
		}
		tprefix = fmt.Sprintf("%s.%d.%d", now_date, h, m)
	}

	if tprefix != "" && rotate {
		new_name := fmt.Sprintf("%s.%s.%s", prefix, tprefix, suffix)
		f.Sync()
		f.Close()
		os.Rename(name, new_name)
		f, _ = os.Create(l.fname)
		l.output = f
		l.raw.SetOutput(f)
		return nil
	}

	return nil
}
