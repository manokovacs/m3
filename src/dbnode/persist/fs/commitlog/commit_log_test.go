// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package commitlog

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/m3db/bitset"
	"github.com/m3db/m3/src/dbnode/persist"
	"github.com/m3db/m3/src/dbnode/persist/fs"
	"github.com/m3db/m3/src/dbnode/ts"
	"github.com/m3db/m3/src/x/context"
	"github.com/m3db/m3/src/x/ident"
	xtime "github.com/m3db/m3/src/x/time"

	mclock "github.com/facebookgo/clock"
	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/require"
	"github.com/uber-go/tally"
)

// readAllSeriesPredicateTest is the same as ReadAllSeriesPredicate except
// it asserts that the ID and the namespace are not nil.
func readAllSeriesPredicateTest() SeriesFilterPredicate {
	return func(id ident.ID, namespace ident.ID) bool {
		if id == nil {
			panic(fmt.Sprintf("series ID passed to series predicate is nil"))
		}

		if namespace == nil {
			panic(fmt.Sprintf("namespace ID passed to series predicate is nil"))
		}

		return true
	}
}

type overrides struct {
	clock            *mclock.Mock
	flushInterval    *time.Duration
	backlogQueueSize *int
	strategy         Strategy
}

var testOpts = NewOptions().
	SetBlockSize(2 * time.Hour).
	SetFlushSize(4096).
	SetFlushInterval(100 * time.Millisecond).
	SetBacklogQueueSize(1024)

func newTestOptions(
	t *testing.T,
	overrides overrides,
) (
	Options,
	tally.TestScope,
) {
	dir, err := ioutil.TempDir("", "foo")
	require.NoError(t, err)

	var c mclock.Clock
	if overrides.clock != nil {
		c = overrides.clock
	} else {
		c = mclock.New()
	}

	scope := tally.NewTestScope("", nil)

	opts := testOpts.
		SetClockOptions(testOpts.ClockOptions().SetNowFn(c.Now)).
		SetInstrumentOptions(testOpts.InstrumentOptions().SetMetricsScope(scope)).
		SetFilesystemOptions(testOpts.FilesystemOptions().SetFilePathPrefix(dir))

	if overrides.flushInterval != nil {
		opts = opts.SetFlushInterval(*overrides.flushInterval)
	}

	if overrides.backlogQueueSize != nil {
		opts = opts.SetBacklogQueueSize(*overrides.backlogQueueSize)
	}

	opts = opts.SetStrategy(overrides.strategy)

	return opts, scope
}

func cleanup(t *testing.T, opts Options) {
	filePathPrefix := opts.FilesystemOptions().FilePathPrefix()
	require.NoError(t, os.RemoveAll(filePathPrefix))
}

type testWrite struct {
	series      ts.Series
	t           time.Time
	v           float64
	u           xtime.Unit
	a           []byte
	expectedErr error
}

func testSeries(
	uniqueIndex uint64,
	id string,
	tags ident.Tags,
	shard uint32,
) ts.Series {
	return ts.Series{
		UniqueIndex: uniqueIndex,
		Namespace:   ident.StringID("testNS"),
		ID:          ident.StringID(id),
		Tags:        tags,
		Shard:       shard,
	}
}

func (w testWrite) assert(
	t *testing.T,
	series ts.Series,
	datapoint ts.Datapoint,
	unit xtime.Unit,
	annotation []byte,
) {
	require.Equal(t, w.series.UniqueIndex, series.UniqueIndex)
	require.True(t, w.series.ID.Equal(series.ID), fmt.Sprintf("write ID '%s' does not match actual ID '%s'", w.series.ID.String(), series.ID.String()))
	require.Equal(t, w.series.Shard, series.Shard)

	// ident.Tags.Equal will compare length
	require.True(t, w.series.Tags.Equal(series.Tags))

	require.True(t, w.t.Equal(datapoint.Timestamp))
	require.Equal(t, w.v, datapoint.Value)
	require.Equal(t, w.u, unit)
	require.Equal(t, w.a, annotation)
}

func snapshotCounterValue(
	scope tally.TestScope,
	counter string,
) (tally.CounterSnapshot, bool) {
	counters := scope.Snapshot().Counters()
	c, ok := counters[tally.KeyForPrefixedStringMap(counter, nil)]
	return c, ok
}

type mockCommitLogWriter struct {
	openFn  func() (persist.CommitLogFile, error)
	writeFn func(ts.Series, ts.Datapoint, xtime.Unit, ts.Annotation) error
	flushFn func(sync bool) error
	closeFn func() error
}

func newMockCommitLogWriter() *mockCommitLogWriter {
	return &mockCommitLogWriter{
		openFn: func() (persist.CommitLogFile, error) {
			return persist.CommitLogFile{}, nil
		},
		writeFn: func(ts.Series, ts.Datapoint, xtime.Unit, ts.Annotation) error {
			return nil
		},
		flushFn: func(sync bool) error {
			return nil
		},
		closeFn: func() error {
			return nil
		},
	}
}

func (w *mockCommitLogWriter) Open() (persist.CommitLogFile, error) {
	return w.openFn()
}

func (w *mockCommitLogWriter) Write(
	series ts.Series,
	datapoint ts.Datapoint,
	unit xtime.Unit,
	annotation ts.Annotation,
) error {
	return w.writeFn(series, datapoint, unit, annotation)
}

func (w *mockCommitLogWriter) Flush(sync bool) error {
	return w.flushFn(sync)
}

func (w *mockCommitLogWriter) Close() error {
	return w.closeFn()
}

func newTestCommitLog(t *testing.T, opts Options) *commitLog {
	commitLogI, err := NewCommitLog(opts)
	require.NoError(t, err)
	commitLog := commitLogI.(*commitLog)
	require.NoError(t, commitLog.Open())

	// Ensure files present
	fsopts := opts.FilesystemOptions()
	files, err := fs.SortedCommitLogFiles(fs.CommitLogsDirPath(fsopts.FilePathPrefix()))
	require.NoError(t, err)
	require.True(t, len(files) == 2)

	return commitLog
}

func writeCommitLogs(
	t *testing.T,
	scope tally.TestScope,
	commitLog CommitLog,
	writes []testWrite,
) *sync.WaitGroup {
	wg := sync.WaitGroup{}

	getAllWrites := func() int {
		result := int64(0)
		success, ok := snapshotCounterValue(scope, "commitlog.writes.success")
		if ok {
			result += success.Value()
		}
		errors, ok := snapshotCounterValue(scope, "commitlog.writes.errors")
		if ok {
			result += errors.Value()
		}
		return int(result)
	}

	ctx := context.NewContext()
	defer ctx.Close()

	preWrites := getAllWrites()

	for i, write := range writes {
		i := i
		write := write

		// Wait for previous writes to enqueue
		for getAllWrites() != preWrites+i {
			time.Sleep(time.Microsecond)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			series := write.series
			datapoint := ts.Datapoint{Timestamp: write.t, Value: write.v}
			err := commitLog.Write(ctx, series, datapoint, write.u, write.a)

			if write.expectedErr != nil {
				if !strings.Contains(fmt.Sprintf("%v", err), fmt.Sprintf("%v", write.expectedErr)) {
					panic(fmt.Sprintf("unexpected error: %v", err))
				}
			} else {
				if err != nil {
					panic(err)
				}
			}
		}()
	}

	// Wait for all writes to enqueue
	for getAllWrites() != preWrites+len(writes) {
		time.Sleep(time.Microsecond)
	}

	return &wg
}

type seriesTestWritesAndReadPosition struct {
	writes       []testWrite
	readPosition int
}

func assertCommitLogWritesByIterating(t *testing.T, l *commitLog, writes []testWrite) {
	iterOpts := IteratorOpts{
		CommitLogOptions:      l.opts,
		FileFilterPredicate:   ReadAllPredicate(),
		SeriesFilterPredicate: readAllSeriesPredicateTest(),
	}
	iter, corruptFiles, err := NewIterator(iterOpts)
	require.NoError(t, err)
	require.Equal(t, 0, len(corruptFiles))
	defer iter.Close()

	// Convert the writes to be in-order, but keyed by series ID because the
	// commitlog reader only guarantees the same order on disk within a
	// given series
	writesBySeries := map[string]seriesTestWritesAndReadPosition{}
	for _, write := range writes {
		seriesWrites := writesBySeries[write.series.ID.String()]
		if seriesWrites.writes == nil {
			seriesWrites.writes = []testWrite{}
		}
		seriesWrites.writes = append(seriesWrites.writes, write)
		writesBySeries[write.series.ID.String()] = seriesWrites
	}

	for iter.Next() {
		series, datapoint, unit, annotation := iter.Current()

		seriesWrites := writesBySeries[series.ID.String()]
		write := seriesWrites.writes[seriesWrites.readPosition]

		write.assert(t, series, datapoint, unit, annotation)

		seriesWrites.readPosition++
		writesBySeries[series.ID.String()] = seriesWrites
	}

	require.NoError(t, iter.Err())
}

func setupCloseOnFail(t *testing.T, l *commitLog) *sync.WaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(1)
	l.commitLogFailFn = func(err error) {
		go func() {
			l.Close()
			wg.Done()
		}()
	}
	return &wg
}

func TestCommitLogWrite(t *testing.T) {
	opts, scope := newTestOptions(t, overrides{
		strategy: StrategyWriteWait,
	})
	defer cleanup(t, opts)

	commitLog := newTestCommitLog(t, opts)

	writes := []testWrite{
		{testSeries(0, "foo.bar", ident.NewTags(ident.StringTag("name1", "val1")), 127), time.Now(), 123.456, xtime.Second, []byte{1, 2, 3}, nil},
		{testSeries(1, "foo.baz", ident.NewTags(ident.StringTag("name2", "val2")), 150), time.Now(), 456.789, xtime.Second, nil, nil},
	}

	// Call write sync
	writeCommitLogs(t, scope, commitLog, writes).Wait()

	// Close the commit log and consequently flush
	require.NoError(t, commitLog.Close())

	// Assert writes occurred by reading the commit log
	assertCommitLogWritesByIterating(t, commitLog, writes)
}

func TestReadCommitLogMissingMetadata(t *testing.T) {
	readConc := 4
	// Make sure we're not leaking goroutines
	defer leaktest.CheckTimeout(t, 10*time.Second)()

	opts, scope := newTestOptions(t, overrides{
		strategy: StrategyWriteWait,
	})
	// Set read concurrency so that the parallel path is definitely tested
	opts.SetReadConcurrency(readConc)
	defer cleanup(t, opts)

	// Replace bitset in writer with one that configurably returns true or false
	// depending on the series
	commitLog := newTestCommitLog(t, opts)
	primary := commitLog.writerState.primary.writer.(*writer)
	secondary := commitLog.writerState.secondary.writer.(*writer)

	bitSet := bitset.NewBitSet(0)

	// Generate fake series, where approximately half will be missing metadata.
	// This works because the commitlog writer uses the bitset to determine if
	// the metadata for a particular series had already been written to disk.
	allSeries := []ts.Series{}
	for i := 0; i < 200; i++ {
		willNotHaveMetadata := !(i%2 == 0)
		allSeries = append(allSeries, testSeries(
			uint64(i),
			"hax",
			ident.NewTags(ident.StringTag("name", "val")),
			uint32(i%100),
		))
		if willNotHaveMetadata {
			bitSet.Set(uint(i))
		}
	}
	primary.seen = bitSet
	secondary.seen = bitSet

	// Generate fake writes for each of the series
	writes := []testWrite{}
	for _, series := range allSeries {
		for i := 0; i < 10; i++ {
			writes = append(writes, testWrite{series, time.Now(), rand.Float64(), xtime.Second, []byte{1, 2, 3}, nil})
		}
	}

	// Call write sync
	writeCommitLogs(t, scope, commitLog, writes).Wait()

	// Close the commit log and consequently flush
	require.NoError(t, commitLog.Close())

	// Make sure we don't panic / deadlock
	iterOpts := IteratorOpts{
		CommitLogOptions:      opts,
		FileFilterPredicate:   ReadAllPredicate(),
		SeriesFilterPredicate: readAllSeriesPredicateTest(),
	}
	iter, corruptFiles, err := NewIterator(iterOpts)
	require.NoError(t, err)
	require.Equal(t, 0, len(corruptFiles))

	for iter.Next() {
		require.NoError(t, iter.Err())
	}
	require.Equal(t, errCommitLogReaderMissingMetadata, iter.Err())
	iter.Close()
	require.NoError(t, commitLog.Close())
}

func TestCommitLogReaderIsNotReusable(t *testing.T) {
	// Make sure we're not leaking goroutines
	defer leaktest.CheckTimeout(t, time.Second)()

	overrideFlushInterval := 10 * time.Millisecond
	opts, scope := newTestOptions(t, overrides{
		strategy:      StrategyWriteWait,
		flushInterval: &overrideFlushInterval,
	})
	defer cleanup(t, opts)

	commitLog := newTestCommitLog(t, opts)

	writes := []testWrite{
		{testSeries(0, "foo.bar", testTags1, 127), time.Now(), 123.456, xtime.Second, []byte{1, 2, 3}, nil},
		{testSeries(1, "foo.baz", testTags2, 150), time.Now(), 456.789, xtime.Second, nil, nil},
	}

	// Call write sync
	writeCommitLogs(t, scope, commitLog, writes).Wait()

	// Close the commit log and consequently flush
	require.NoError(t, commitLog.Close())

	// Assert writes occurred by reading the commit log
	assertCommitLogWritesByIterating(t, commitLog, writes)

	// Assert commitlog file exists and retrieve path
	fsopts := opts.FilesystemOptions()
	files, err := fs.SortedCommitLogFiles(fs.CommitLogsDirPath(fsopts.FilePathPrefix()))
	require.NoError(t, err)
	require.Equal(t, 2, len(files))

	// Assert commitlog cannot be opened more than once
	reader := newCommitLogReader(opts, readAllSeriesPredicateTest())
	_, err = reader.Open(files[0])
	require.NoError(t, err)
	reader.Close()
	_, err = reader.Open(files[0])
	require.Equal(t, errCommitLogReaderIsNotReusable, err)
}

func TestCommitLogIteratorUsesPredicateFilterForNonCorruptFiles(t *testing.T) {
	clock := mclock.NewMock()
	opts, scope := newTestOptions(t, overrides{
		clock:    clock,
		strategy: StrategyWriteWait,
	})

	start := clock.Now()

	// Writes spaced apart by block size.
	writes := []testWrite{
		{testSeries(0, "foo.bar", testTags1, 127), start, 123.456, xtime.Millisecond, nil, nil},
		{testSeries(1, "foo.baz", testTags2, 150), start.Add(1 * time.Second), 456.789, xtime.Millisecond, nil, nil},
		{testSeries(2, "foo.qux", testTags3, 291), start.Add(2 * time.Second), 789.123, xtime.Millisecond, nil, nil},
	}
	defer cleanup(t, opts)

	commitLog := newTestCommitLog(t, opts)

	// Write, making sure that the clock is set properly for each write.
	for _, write := range writes {
		// Modify the time to make sure we're generating commitlog files with different
		// start times.
		clock.Add(write.t.Sub(clock.Now()))
		// Rotate frequently to ensure we're generating multiple files.
		_, err := commitLog.RotateLogs()
		require.NoError(t, err)
		writeCommitLogs(t, scope, commitLog, []testWrite{write})
	}

	// Close the commit log and consequently flush.
	require.NoError(t, commitLog.Close())

	// Make sure multiple commitlog files were generated.
	fsopts := opts.FilesystemOptions()
	files, err := fs.SortedCommitLogFiles(fs.CommitLogsDirPath(fsopts.FilePathPrefix()))
	require.NoError(t, err)
	require.Equal(t, 5, len(files))

	// This predicate should eliminate the first commitlog file.
	commitLogPredicate := func(f FileFilterInfo) bool {
		require.False(t, f.IsCorrupt)
		return f.File.Index > 0
	}

	// Assert that the commitlog iterator honors the predicate and only uses
	// 2 of the 3 files.
	iterOpts := IteratorOpts{
		CommitLogOptions:      opts,
		FileFilterPredicate:   commitLogPredicate,
		SeriesFilterPredicate: readAllSeriesPredicateTest(),
	}
	iter, corruptFiles, err := NewIterator(iterOpts)
	require.NoError(t, err)
	require.True(t, len(corruptFiles) <= 1)

	iterStruct := iter.(*iterator)
	require.True(t, len(iterStruct.files) >= 4)
}

func TestCommitLogIteratorUsesPredicateFilterForCorruptFiles(t *testing.T) {
	clock := mclock.NewMock()
	opts, _ := newTestOptions(t, overrides{
		clock:    clock,
		strategy: StrategyWriteWait,
	})
	defer cleanup(t, opts)

	commitLog := newTestCommitLog(t, opts)
	// Close the commit log and consequently flush.
	require.NoError(t, commitLog.Close())

	// Make sure a valid commitlog was created.
	fsopts := opts.FilesystemOptions()
	files, err := fs.SortedCommitLogFiles(fs.CommitLogsDirPath(fsopts.FilePathPrefix()))
	require.NoError(t, err)
	require.Equal(t, 2, len(files))

	// Write out a corrupt commitlog file.
	nextCommitlogFilePath, _, err := NextFile(opts)
	require.NoError(t, err)
	err = ioutil.WriteFile(
		nextCommitlogFilePath, []byte("not-a-valid-commitlog-file"), os.FileMode(0666))
	require.NoError(t, err)

	// Make sure the corrupt file is visibile.
	files, err = fs.SortedCommitLogFiles(fs.CommitLogsDirPath(fsopts.FilePathPrefix()))
	require.NoError(t, err)
	require.Equal(t, 3, len(files))

	// Assert that the corrupt file is returned from the iterator.
	iterOpts := IteratorOpts{
		CommitLogOptions:      opts,
		FileFilterPredicate:   ReadAllPredicate(),
		SeriesFilterPredicate: readAllSeriesPredicateTest(),
	}
	iter, corruptFiles, err := NewIterator(iterOpts)
	require.NoError(t, err)
	require.Equal(t, 1, len(corruptFiles))

	iterStruct := iter.(*iterator)
	require.Equal(t, 2, len(iterStruct.files))

	// Assert that the iterator ignores the corrupt file given an appropriate predicate.
	ignoreCorruptPredicate := func(f FileFilterInfo) bool {
		return !f.IsCorrupt
	}

	iterOpts = IteratorOpts{
		CommitLogOptions:      opts,
		FileFilterPredicate:   ignoreCorruptPredicate,
		SeriesFilterPredicate: readAllSeriesPredicateTest(),
	}
	iter, corruptFiles, err = NewIterator(iterOpts)
	require.NoError(t, err)
	require.Equal(t, 0, len(corruptFiles))

	iterStruct = iter.(*iterator)
	require.Equal(t, 2, len(iterStruct.files))
}

func TestCommitLogWriteBehind(t *testing.T) {
	opts, scope := newTestOptions(t, overrides{
		strategy: StrategyWriteBehind,
	})
	defer cleanup(t, opts)

	commitLog := newTestCommitLog(t, opts)

	writes := []testWrite{
		{testSeries(0, "foo.bar", testTags1, 127), time.Now(), 123.456, xtime.Millisecond, nil, nil},
		{testSeries(1, "foo.baz", testTags2, 150), time.Now(), 456.789, xtime.Millisecond, nil, nil},
		{testSeries(2, "foo.qux", testTags3, 291), time.Now(), 789.123, xtime.Millisecond, nil, nil},
	}

	// Call write behind
	writeCommitLogs(t, scope, commitLog, writes)

	// Close the commit log and consequently flush
	require.NoError(t, commitLog.Close())

	// Assert writes occurred by reading the commit log
	assertCommitLogWritesByIterating(t, commitLog, writes)
}

func TestCommitLogWriteErrorOnClosed(t *testing.T) {
	opts, _ := newTestOptions(t, overrides{})
	defer cleanup(t, opts)

	commitLog := newTestCommitLog(t, opts)
	require.NoError(t, commitLog.Close())

	series := testSeries(0, "foo.bar", testTags1, 127)
	datapoint := ts.Datapoint{Timestamp: time.Now(), Value: 123.456}

	ctx := context.NewContext()
	defer ctx.Close()

	err := commitLog.Write(ctx, series, datapoint, xtime.Millisecond, nil)
	require.Error(t, err)
	require.Equal(t, errCommitLogClosed, err)
}

func TestCommitLogWriteErrorOnFull(t *testing.T) {
	// Set backlog of size one and don't automatically flush.
	backlogQueueSize := 1
	flushInterval := time.Duration(0)
	opts, _ := newTestOptions(t, overrides{
		backlogQueueSize: &backlogQueueSize,
		flushInterval:    &flushInterval,
		strategy:         StrategyWriteBehind,
	})
	defer cleanup(t, opts)

	commitLog := newTestCommitLog(t, opts)

	// Test filling queue
	var writes []testWrite
	series := testSeries(0, "foo.bar", testTags1, 127)
	dp := ts.Datapoint{Timestamp: time.Now(), Value: 123.456}
	unit := xtime.Millisecond

	ctx := context.NewContext()
	defer ctx.Close()

	for {
		if err := commitLog.Write(ctx, series, dp, unit, nil); err != nil {
			// Ensure queue full error.
			require.Equal(t, ErrCommitLogQueueFull, err)
			require.Equal(t, int64(backlogQueueSize), commitLog.QueueLength())
			break
		}
		writes = append(writes, testWrite{series, dp.Timestamp, dp.Value, unit, nil, nil})

		// Increment timestamp and value for next write.
		dp.Timestamp = dp.Timestamp.Add(time.Second)
		dp.Value += 1.0
	}

	// Close and consequently flush.
	require.NoError(t, commitLog.Close())

	// Assert write flushed by reading the commit log.
	assertCommitLogWritesByIterating(t, commitLog, writes)
}

func TestCommitLogQueueLength(t *testing.T) {
	// Set backlog of size one and don't automatically flush.
	backlogQueueSize := 10
	flushInterval := time.Duration(0)
	opts, _ := newTestOptions(t, overrides{
		backlogQueueSize: &backlogQueueSize,
		flushInterval:    &flushInterval,
		strategy:         StrategyWriteBehind,
	})
	defer cleanup(t, opts)

	commitLog := newTestCommitLog(t, opts)
	defer commitLog.Close()

	var (
		series = testSeries(0, "foo.bar", testTags1, 127)
		dp     = ts.Datapoint{Timestamp: time.Now(), Value: 123.456}
		unit   = xtime.Millisecond
		ctx    = context.NewContext()
	)
	defer ctx.Close()

	for i := 0; ; i++ {
		// Write in a loop and check the queue length until the queue is full.
		require.Equal(t, int64(i), commitLog.QueueLength())
		if err := commitLog.Write(ctx, series, dp, unit, nil); err != nil {
			require.Equal(t, ErrCommitLogQueueFull, err)
			break
		}

		// Increment timestamp and value for next write.
		dp.Timestamp = dp.Timestamp.Add(time.Second)
		dp.Value += 1.0
	}
}

func TestCommitLogFailOnWriteError(t *testing.T) {
	opts, scope := newTestOptions(t, overrides{
		strategy: StrategyWriteBehind,
	})
	defer cleanup(t, opts)

	commitLogI, err := NewCommitLog(opts)
	require.NoError(t, err)
	commitLog := commitLogI.(*commitLog)
	writer := newMockCommitLogWriter()

	writer.writeFn = func(ts.Series, ts.Datapoint, xtime.Unit, ts.Annotation) error {
		return fmt.Errorf("an error")
	}

	writer.openFn = func() (persist.CommitLogFile, error) {
		return persist.CommitLogFile{}, nil
	}

	writer.flushFn = func(bool) error {
		commitLog.writerState.primary.onFlush(nil)
		return nil
	}

	commitLog.newCommitLogWriterFn = func(
		_ flushFn,
		_ Options,
	) commitLogWriter {
		return writer
	}

	require.NoError(t, commitLog.Open())

	wg := setupCloseOnFail(t, commitLog)

	writes := []testWrite{
		{testSeries(0, "foo.bar", testTags1, 127), time.Now(), 123.456, xtime.Millisecond, nil, nil},
	}

	writeCommitLogs(t, scope, commitLog, writes)

	wg.Wait()

	// Check stats
	errors, ok := snapshotCounterValue(scope, "commitlog.writes.errors")
	require.True(t, ok)
	require.Equal(t, int64(1), errors.Value())
}

func TestCommitLogFailOnOpenError(t *testing.T) {
	opts, scope := newTestOptions(t, overrides{
		strategy: StrategyWriteBehind,
	})
	defer cleanup(t, opts)

	commitLogI, err := NewCommitLog(opts)
	require.NoError(t, err)
	commitLog := commitLogI.(*commitLog)
	writer := newMockCommitLogWriter()

	var opens int64
	writer.openFn = func() (persist.CommitLogFile, error) {
		if atomic.AddInt64(&opens, 1) >= 3 {
			return persist.CommitLogFile{}, fmt.Errorf("an error")
		}
		return persist.CommitLogFile{}, nil
	}

	writer.flushFn = func(bool) error {
		commitLog.writerState.primary.onFlush(nil)
		return nil
	}

	commitLog.newCommitLogWriterFn = func(
		_ flushFn,
		_ Options,
	) commitLogWriter {
		return writer
	}

	require.NoError(t, commitLog.Open())

	wg := setupCloseOnFail(t, commitLog)

	writes := []testWrite{
		{testSeries(0, "foo.bar", testTags1, 127), time.Now(), 123.456, xtime.Millisecond, nil, nil},
	}

	writeCommitLogs(t, scope, commitLog, writes)

	// Rotate the commitlog so that it requires a new open.
	commitLog.RotateLogs()

	wg.Wait()
	// Secondary writer open is async so wait for it to complete before asserting
	// that it failed.
	commitLog.waitForSecondaryWriterAsyncResetComplete()

	// Check stats
	errors, ok := snapshotCounterValue(scope, "commitlog.writes.errors")
	require.True(t, ok)
	require.Equal(t, int64(1), errors.Value())

	openErrors, ok := snapshotCounterValue(scope, "commitlog.writes.open-errors")
	require.True(t, ok)
	require.Equal(t, int64(1), openErrors.Value())
}

func TestCommitLogFailOnFlushError(t *testing.T) {
	opts, scope := newTestOptions(t, overrides{
		strategy: StrategyWriteBehind,
	})
	defer cleanup(t, opts)

	commitLogI, err := NewCommitLog(opts)
	require.NoError(t, err)
	commitLog := commitLogI.(*commitLog)
	writer := newMockCommitLogWriter()

	var flushes int64
	writer.flushFn = func(bool) error {
		if atomic.AddInt64(&flushes, 1) >= 2 {
			commitLog.writerState.primary.onFlush(fmt.Errorf("an error"))
		} else {
			commitLog.writerState.primary.onFlush(nil)
		}
		return nil
	}

	commitLog.newCommitLogWriterFn = func(
		_ flushFn,
		_ Options,
	) commitLogWriter {
		return writer
	}

	require.NoError(t, commitLog.Open())

	wg := setupCloseOnFail(t, commitLog)

	writes := []testWrite{
		{testSeries(0, "foo.bar", testTags1, 127), time.Now(), 123.456, xtime.Millisecond, nil, nil},
	}

	writeCommitLogs(t, scope, commitLog, writes)

	wg.Wait()

	// Check stats
	errors, ok := snapshotCounterValue(scope, "commitlog.writes.errors")
	require.True(t, ok)
	require.Equal(t, int64(2), errors.Value())

	flushErrors, ok := snapshotCounterValue(scope, "commitlog.writes.flush-errors")
	require.True(t, ok)
	require.Equal(t, int64(2), flushErrors.Value())
}

func TestCommitLogActiveLogs(t *testing.T) {
	opts, _ := newTestOptions(t, overrides{
		strategy: StrategyWriteBehind,
	})
	defer cleanup(t, opts)

	commitLog := newTestCommitLog(t, opts)

	writer := newMockCommitLogWriter()
	writer.flushFn = func(bool) error {
		return nil
	}
	commitLog.newCommitLogWriterFn = func(
		_ flushFn,
		_ Options,
	) commitLogWriter {
		return writer
	}

	logs, err := commitLog.ActiveLogs()
	require.NoError(t, err)
	require.Equal(t, 2, len(logs))

	// Close the commit log and consequently flush
	require.NoError(t, commitLog.Close())
	_, err = commitLog.ActiveLogs()
	require.Error(t, err)
}

func TestCommitLogRotateLogs(t *testing.T) {
	var (
		clock       = mclock.NewMock()
		opts, scope = newTestOptions(t, overrides{
			clock:    clock,
			strategy: StrategyWriteWait,
		})
	)
	defer cleanup(t, opts)

	var (
		commitLog = newTestCommitLog(t, opts)
		start     = clock.Now()
	)

	// Writes spaced such that they should appear within the same commitlog block.
	writes := []testWrite{
		{testSeries(0, "foo.bar", testTags1, 127), start, 123.456, xtime.Millisecond, nil, nil},
		{testSeries(1, "foo.baz", testTags2, 150), start.Add(1 * time.Second), 456.789, xtime.Millisecond, nil, nil},
		{testSeries(2, "foo.qux", testTags3, 291), start.Add(2 * time.Second), 789.123, xtime.Millisecond, nil, nil},
	}

	for i, write := range writes {
		// Set clock to align with the write.
		clock.Add(write.t.Sub(clock.Now()))

		// Write entry.
		writeCommitLogs(t, scope, commitLog, []testWrite{write})

		file, err := commitLog.RotateLogs()
		require.NoError(t, err)
		require.Equal(t, file.Index, int64(i+1))
		require.Contains(t, file.FilePath, "commitlog-0")
	}

	// Secondary writer open is async so wait for it to complete so that its safe to assert
	// on the number of files that should be on disk otherwise test will flake depending
	// on whether or not the async open completed in time.
	commitLog.waitForSecondaryWriterAsyncResetComplete()

	// Ensure files present for each call to RotateLogs().
	fsopts := opts.FilesystemOptions()
	files, err := fs.SortedCommitLogFiles(fs.CommitLogsDirPath(fsopts.FilePathPrefix()))
	require.NoError(t, err)
	require.Equal(t, len(writes)+2, len(files)) // +2 to account for the initial files.

	// Close and consequently flush.
	require.NoError(t, commitLog.Close())

	// Assert write flushed by reading the commit log.
	assertCommitLogWritesByIterating(t, commitLog, writes)
}

var (
	testTag1 = ident.StringTag("name1", "val1")
	testTag2 = ident.StringTag("name2", "val2")
	testTag3 = ident.StringTag("name3", "val3")

	testTags1 = ident.NewTags(testTag1)
	testTags2 = ident.NewTags(testTag2)
	testTags3 = ident.NewTags(testTag3)
)

func TestCommitLogBatchWriteDoesNotAddErroredOrSkippedSeries(t *testing.T) {
	opts, scope := newTestOptions(t, overrides{
		strategy: StrategyWriteWait,
	})

	defer cleanup(t, opts)
	commitLog := newTestCommitLog(t, opts)
	finalized := 0
	finalizeFn := func(_ ts.WriteBatch) {
		finalized++
	}

	writes := ts.NewWriteBatch(4, ident.StringID("ns"), finalizeFn)

	clock := mclock.NewMock()
	alignedStart := clock.Now().Truncate(time.Hour)
	for i := 0; i < 4; i++ {
		tt := alignedStart.Add(time.Minute * time.Duration(i))
		writes.Add(i, ident.StringID(fmt.Sprint(i)), tt, float64(i)*10.5, xtime.Second, nil)
	}

	writes.SetSkipWrite(0)
	writes.SetOutcome(1, testSeries(1, "foo.bar", testTags1, 127), nil)
	writes.SetOutcome(2, testSeries(2, "err.err", testTags2, 255), errors.New("oops"))
	writes.SetOutcome(3, testSeries(3, "biz.qux", testTags3, 511), nil)

	// Call write batch sync
	wg := sync.WaitGroup{}

	getAllWrites := func() int {
		result := int64(0)
		success, ok := snapshotCounterValue(scope, "commitlog.writes.success")
		if ok {
			result += success.Value()
		}
		errors, ok := snapshotCounterValue(scope, "commitlog.writes.errors")
		if ok {
			result += errors.Value()
		}
		return int(result)
	}

	ctx := context.NewContext()
	defer ctx.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := commitLog.WriteBatch(ctx, writes)
		require.NoError(t, err)
	}()

	// Wait for all writes to enqueue
	for getAllWrites() != 2 {
		time.Sleep(time.Microsecond)
	}

	wg.Wait()

	// Close the commit log and consequently flush
	require.NoError(t, commitLog.Close())

	// Assert writes occurred by reading the commit log
	expected := []testWrite{
		{testSeries(1, "foo.bar", testTags1, 127), alignedStart.Add(time.Minute), 10.5, xtime.Second, nil, nil},
		{testSeries(3, "biz.qux", testTags3, 511), alignedStart.Add(time.Minute * 3), 31.5, xtime.Second, nil, nil},
	}

	assertCommitLogWritesByIterating(t, commitLog, expected)
	require.Equal(t, 1, finalized)
}
