// Copyright (c) 2019 Uber Technologies, Inc.
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

package storage

import (
	"errors"
	"testing"

	"github.com/m3db/m3/src/dbnode/namespace"
	"github.com/m3db/m3/src/dbnode/storage/series"
	"github.com/m3db/m3/src/dbnode/x/xio"
	"github.com/m3db/m3/src/x/context"
	"github.com/m3db/m3/src/x/ident"
	xtime "github.com/m3db/m3/src/x/time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dirtyData struct {
	id    ident.ID
	start xtime.UnixNano
}

func TestRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := NewMockdatabaseShard(ctrl)
	retriever := series.NewMockQueryableBlockRetriever(ctrl)
	version := 0
	ctx := context.NewContext()
	nsCtx := namespace.Context{}
	fetchedBlocks := []xio.BlockReader{xio.BlockReader{}}
	retriever.EXPECT().RetrievableBlockColdVersion(gomock.Any()).Return(version).AnyTimes()

	dirtySeries := newDirtySeriesMap(dirtySeriesMapOptions{})
	dirtySeriesToWrite := make(map[xtime.UnixNano]*idList)

	data := []dirtyData{
		dirtyData{start: 0, id: ident.StringID("id0")},
		dirtyData{start: 0, id: ident.StringID("id1")},
		dirtyData{start: 1, id: ident.StringID("id2")},
		dirtyData{start: 1, id: ident.StringID("id3")},
		dirtyData{start: 1, id: ident.StringID("id4")},
		dirtyData{start: 2, id: ident.StringID("id5")},
		dirtyData{start: 3, id: ident.StringID("id6")},
		dirtyData{start: 3, id: ident.StringID("id7")},
		dirtyData{start: 4, id: ident.StringID("id8")},
	}

	// Populate bookkeeping data structures with above test data.
	for _, d := range data {
		addDirtySeries(dirtySeries, dirtySeriesToWrite, d.id, d.start)
		shard.EXPECT().
			FetchBlocksForColdFlush(gomock.Any(), d.id, d.start.ToTime(), version+1, nsCtx).
			Return(fetchedBlocks, nil)
	}

	mergeWith := newFSMergeWithMem(shard, retriever, dirtySeries, dirtySeriesToWrite)

	for _, d := range data {
		require.True(t, dirtySeries.Contains(idAndBlockStart{blockStart: d.start, id: d.id}))
		beforeLen := dirtySeriesToWrite[d.start].Len()
		res, exists, err := mergeWith.Read(ctx, d.id, d.start, nsCtx)
		require.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, fetchedBlocks, res)
		// Assert that the Read call removes the element from the "to write"
		// list.
		assert.Equal(t, beforeLen-1, dirtySeriesToWrite[d.start].Len())
	}

	// Test Read with non-existent dirty block/series.
	res, exists, err := mergeWith.Read(ctx, ident.StringID("not-present"), 10, nsCtx)
	assert.Nil(t, res)
	assert.False(t, exists)
	assert.NoError(t, err)

	// Test Read with error on fetch.
	badFetchID := ident.StringID("bad-fetch")
	addDirtySeries(dirtySeries, dirtySeriesToWrite, badFetchID, 11)
	shard.EXPECT().
		FetchBlocksForColdFlush(gomock.Any(), badFetchID, gomock.Any(), version+1, nsCtx).
		Return(nil, errors.New("fetch error"))
	res, exists, err = mergeWith.Read(ctx, badFetchID, 11, nsCtx)
	assert.Nil(t, res)
	assert.False(t, exists)
	assert.Error(t, err)

	// Test Read with no data on fetch.
	emptyDataID := ident.StringID("empty-data")
	addDirtySeries(dirtySeries, dirtySeriesToWrite, emptyDataID, 12)
	shard.EXPECT().
		FetchBlocksForColdFlush(gomock.Any(), emptyDataID, gomock.Any(), version+1, nsCtx).
		Return(nil, nil)
	res, exists, err = mergeWith.Read(ctx, emptyDataID, 12, nsCtx)
	assert.Nil(t, res)
	assert.False(t, exists)
	assert.NoError(t, err)
}

func TestForEachRemaining(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := NewMockdatabaseShard(ctrl)
	retriever := series.NewMockQueryableBlockRetriever(ctrl)
	version := 0
	ctx := context.NewContext()
	nsCtx := namespace.Context{}
	fetchedBlocks := []xio.BlockReader{xio.BlockReader{}}
	retriever.EXPECT().RetrievableBlockColdVersion(gomock.Any()).Return(version).AnyTimes()

	dirtySeries := newDirtySeriesMap(dirtySeriesMapOptions{})
	dirtySeriesToWrite := make(map[xtime.UnixNano]*idList)

	id0 := ident.StringID("id0")
	id1 := ident.StringID("id1")
	id2 := ident.StringID("id2")
	id3 := ident.StringID("id3")
	id4 := ident.StringID("id4")
	id5 := ident.StringID("id5")
	id6 := ident.StringID("id6")
	id7 := ident.StringID("id7")
	id8 := ident.StringID("id8")
	data := []dirtyData{
		dirtyData{start: 0, id: id0},
		dirtyData{start: 0, id: id1},
		dirtyData{start: 1, id: id2},
		dirtyData{start: 1, id: id3},
		dirtyData{start: 1, id: id4},
		dirtyData{start: 2, id: id5},
		dirtyData{start: 3, id: id6},
		dirtyData{start: 3, id: id7},
		dirtyData{start: 4, id: id8},
	}

	// Populate bookkeeping data structures with above test data.
	for _, d := range data {
		addDirtySeries(dirtySeries, dirtySeriesToWrite, d.id, d.start)
	}

	mergeWith := newFSMergeWithMem(shard, retriever, dirtySeries, dirtySeriesToWrite)

	var forEachCalls []ident.ID
	shard.EXPECT().TagsFromSeriesID(gomock.Any()).Return(ident.Tags{}, true, nil).Times(2)
	shard.EXPECT().
		FetchBlocksForColdFlush(gomock.Any(), id0, xtime.UnixNano(0).ToTime(), version+1, gomock.Any()).
		Return(fetchedBlocks, nil)
	shard.EXPECT().
		FetchBlocksForColdFlush(gomock.Any(), id1, xtime.UnixNano(0).ToTime(), version+1, gomock.Any()).
		Return(fetchedBlocks, nil)
	mergeWith.ForEachRemaining(ctx, 0, func(seriesID ident.ID, tags ident.Tags, data []xio.BlockReader) error {
		forEachCalls = append(forEachCalls, seriesID)
		return nil
	}, nsCtx)
	require.Len(t, forEachCalls, 2)
	assert.Equal(t, id0, forEachCalls[0])
	assert.Equal(t, id1, forEachCalls[1])

	// Reset expected calls.
	forEachCalls = forEachCalls[:0]
	// Read id3 at block start 1, so id2 and id4 should be remaining for block
	// start 1.
	shard.EXPECT().
		FetchBlocksForColdFlush(gomock.Any(), id3, xtime.UnixNano(1).ToTime(), version+1, nsCtx).
		Return(fetchedBlocks, nil)
	res, exists, err := mergeWith.Read(ctx, id3, 1, nsCtx)
	require.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, fetchedBlocks, res)
	shard.EXPECT().TagsFromSeriesID(gomock.Any()).Return(ident.Tags{}, true, nil).Times(2)
	shard.EXPECT().
		FetchBlocksForColdFlush(gomock.Any(), id2, xtime.UnixNano(1).ToTime(), version+1, gomock.Any()).
		Return(fetchedBlocks, nil)
	shard.EXPECT().
		FetchBlocksForColdFlush(gomock.Any(), id4, xtime.UnixNano(1).ToTime(), version+1, gomock.Any()).
		Return(fetchedBlocks, nil)
	err = mergeWith.ForEachRemaining(ctx, 1, func(seriesID ident.ID, tags ident.Tags, data []xio.BlockReader) error {
		forEachCalls = append(forEachCalls, seriesID)
		return nil
	}, nsCtx)
	require.NoError(t, err)
	require.Len(t, forEachCalls, 2)
	assert.Equal(t, id2, forEachCalls[0])
	assert.Equal(t, id4, forEachCalls[1])

	// Test call with error getting tags.
	shard.EXPECT().
		TagsFromSeriesID(gomock.Any()).Return(ident.Tags{}, false, errors.New("bad-tags"))
	shard.EXPECT().
		FetchBlocksForColdFlush(gomock.Any(), id8, xtime.UnixNano(4).ToTime(), version+1, gomock.Any()).
		Return(fetchedBlocks, nil)
	err = mergeWith.ForEachRemaining(ctx, 4, func(seriesID ident.ID, tags ident.Tags, data []xio.BlockReader) error {
		// This function won't be called with the above error.
		return errors.New("unreachable")
	}, nsCtx)
	assert.Error(t, err)

	// Test call with bad function execution.
	shard.EXPECT().
		TagsFromSeriesID(gomock.Any()).Return(ident.Tags{}, true, nil)
	err = mergeWith.ForEachRemaining(ctx, 4, func(seriesID ident.ID, tags ident.Tags, data []xio.BlockReader) error {
		return errors.New("bad")
	}, nsCtx)
	assert.Error(t, err)
}

func addDirtySeries(
	dirtySeries *dirtySeriesMap,
	dirtySeriesToWrite map[xtime.UnixNano]*idList,
	id ident.ID,
	start xtime.UnixNano,
) {
	seriesList := dirtySeriesToWrite[start]
	if seriesList == nil {
		seriesList = newIDList(nil)
		dirtySeriesToWrite[start] = seriesList
	}
	element := seriesList.PushBack(id)

	dirtySeries.Set(idAndBlockStart{blockStart: start, id: id}, element)
}
