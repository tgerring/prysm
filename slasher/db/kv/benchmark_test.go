package kv

import (
	"context"
	"testing"

	basetypes "github.com/farazdagi/prysm-shared-types"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
	"github.com/prysmaticlabs/prysm/slasher/detection/attestations/types"
)

const (
	benchmarkValidator = 300000
)

func BenchmarkStore_SaveEpochSpans(b *testing.B) {
	ctx := context.Background()
	sigBytes := [2]byte{}
	db := setupDB(b)
	es := &types.EpochStore{}

	es, err := es.SetValidatorSpan(benchmarkValidator, types.Span{MinSpan: 1, MaxSpan: 2, SigBytes: sigBytes, HasAttested: true})
	assert.NoError(b, err)
	for i := 0; i < benchmarkValidator; i++ {
		es, err = es.SetValidatorSpan(uint64(i), types.Span{MinSpan: 1, MaxSpan: 2, SigBytes: sigBytes, HasAttested: true})
		assert.NoError(b, err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := db.SaveEpochSpans(ctx, basetypes.Epoch(i%54000), es, false)
		require.NoError(b, err, "Save validator span map failed")
	}
}

func BenchmarkStore_EpochSpans(b *testing.B) {
	db := setupDB(b)
	ctx := context.Background()
	sigBytes := [2]byte{}
	es := &types.EpochStore{}
	es, err := es.SetValidatorSpan(benchmarkValidator, types.Span{MinSpan: 1, MaxSpan: 2, SigBytes: sigBytes, HasAttested: true})
	assert.NoError(b, err)
	for i := 0; i < benchmarkValidator; i++ {
		es, err = es.SetValidatorSpan(uint64(i), types.Span{MinSpan: 1, MaxSpan: 2, SigBytes: sigBytes, HasAttested: true})
		assert.NoError(b, err)
	}
	b.Log(len(es.Bytes()))
	for i := 0; i < 200; i++ {
		err := db.SaveEpochSpans(ctx, basetypes.Epoch(i), es, false)
		require.NoError(b, err, "Save validator span map failed")
	}
	b.Log(db.db.Info())
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := db.EpochSpans(ctx, basetypes.Epoch(i%200), false)
		require.NoError(b, err, "Read validator span map failed")
	}
}

func BenchmarkStore_GetValidatorSpan(b *testing.B) {
	sigBytes := [2]byte{}
	es := &types.EpochStore{}
	es, err := es.SetValidatorSpan(benchmarkValidator, types.Span{MinSpan: 1, MaxSpan: 2, SigBytes: sigBytes, HasAttested: true})
	assert.NoError(b, err)
	for i := 0; i < benchmarkValidator; i++ {
		es, err = es.SetValidatorSpan(uint64(i), types.Span{MinSpan: uint16(i), MaxSpan: uint16(benchmarkValidator - i), SigBytes: sigBytes, HasAttested: true})
		assert.NoError(b, err)
	}
	b.Log(len(es.Bytes()))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := es.GetValidatorSpan(uint64(i % benchmarkValidator))
		require.NoError(b, err, "Read validator span map failed")
	}
}

func BenchmarkStore_SetValidatorSpan(b *testing.B) {
	sigBytes := [2]byte{}
	var err error
	es := &types.EpochStore{}
	es, err = es.SetValidatorSpan(benchmarkValidator, types.Span{MinSpan: 1, MaxSpan: 2, SigBytes: sigBytes, HasAttested: true})
	assert.NoError(b, err)

	for i := 0; i < benchmarkValidator; i++ {
		es, err = es.SetValidatorSpan(uint64(i), types.Span{MinSpan: uint16(i), MaxSpan: uint16(benchmarkValidator - i), SigBytes: sigBytes, HasAttested: true})
		assert.NoError(b, err)
	}
	b.Log(len(es.Bytes()))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		es, err = es.SetValidatorSpan(uint64(i%benchmarkValidator), types.Span{MinSpan: uint16(i), MaxSpan: uint16(benchmarkValidator - i), SigBytes: sigBytes, HasAttested: true})
		require.NoError(b, err, "Read validator span map failed")
	}
}
