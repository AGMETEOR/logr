package target_test

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/wiggin77/logr"
	"github.com/wiggin77/logr/format"
	"github.com/wiggin77/logr/target"
	"github.com/wiggin77/logr/test"
)

func ExampleWriter() {
	lgr := &logr.Logr{}
	buf := &test.Buffer{}
	filter := &logr.StdFilter{Lvl: logr.Warn, Stacktrace: logr.Error}
	formatter := &format.Plain{Delim: " | "}
	t := target.NewWriterTarget(filter, formatter, buf, 1000)
	lgr.AddTarget(t)

	logger := lgr.NewLogger().WithField("name", "wiggin")

	logger.Errorf("the erroneous data is %s", test.StringRnd(10))
	logger.Warnf("strange data: %s", test.StringRnd(5))
	logger.Debug("XXX")
	logger.Trace("XXX")

	err := lgr.Shutdown()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	output := buf.String()
	fmt.Println(output)
}

func TestWriterPlain(t *testing.T) {
	plain := &format.Plain{Delim: " | "}
	writer(t, plain)
}

func TestWriterJSON(t *testing.T) {
	json := &format.JSON{Indent: "  "}
	writer(t, json)
}

func writer(t *testing.T, formatter logr.Formatter) {
	lgr := logr.Logr{}
	buf := &test.Buffer{}
	filter := &logr.StdFilter{Lvl: logr.Warn, Stacktrace: logr.Error}
	target := target.NewWriterTarget(filter, formatter, buf, 1000)
	lgr.AddTarget(target)

	wg := sync.WaitGroup{}
	var id int32

	runner := func(loops int) {
		defer wg.Done()
		tid := atomic.AddInt32(&id, 1)
		logger := lgr.NewLogger().WithFields(logr.Fields{"id": tid, "rnd": rand.Intn(100)})

		for i := 0; i < loops; i++ {
			logger.Debug("XXX")
			logger.Trace("XXX")
			logger.Errorf("count:%d -- the erroneous data is %s", i, test.StringRnd(10))
			logger.Warnf("strange data: %s", test.StringRnd(5))
			runtime.Gosched()
		}
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go runner(50)
	}
	wg.Wait()

	err := lgr.Shutdown()
	if err != nil {
		t.Error(err)
	}

	output := buf.String()
	fmt.Println(output)

	if !strings.Contains(output, "strange data") {
		t.Errorf("missing warnings")
	}

	if strings.Contains(output, "XXX") {
		t.Errorf("wrong level(s) enabled")
	}
}