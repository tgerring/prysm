package kv

import (
	basetypes "github.com/farazdagi/prysm-shared-types"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/slasher/db/types"
)

const (
	latestEpochKey = "LATEST_EPOCH_DETECTED"
	chainHeadKey   = "CHAIN_HEAD"
)

var (
	indexedAttestationsRootsByTargetBucket = []byte("indexed-attestations-roots-by-target")
	indexedAttestationsBucket              = []byte("indexed-attestations")
	// Slasher-related buckets.
	historicIndexedAttestationsBucket = []byte("historic-indexed-attestations-bucket")
	historicBlockHeadersBucket        = []byte("historic-block-headers-bucket")
	highestAttestationBucket          = []byte("highest-attestation-bucket")
	slashingBucket                    = []byte("slashing-bucket")
	chainDataBucket                   = []byte("chain-data-bucket")
	compressedIdxAttsBucket           = []byte("compressed-idx-atts-bucket")
	validatorsPublicKeysBucket        = []byte("validators-public-keys-bucket")
	// In order to quickly detect surround and surrounded attestations we need to store
	// the min and max span for each validator for each epoch.
	// see https://github.com/protolambda/eth2-surround/blob/master/README.md#min-max-surround
	validatorsMinMaxSpanBucket    = []byte("validators-min-max-span-bucket")
	validatorsMinMaxSpanBucketNew = []byte("validators-min-max-span-bucket-new")
)

func encodeSlotValidatorID(slot basetypes.Slot, validatorID uint64) []byte {
	return append(bytesutil.Bytes8(uint64(slot)), bytesutil.Bytes8(validatorID)...)
}

func encodeSlotValidatorIDSig(slot basetypes.Slot, validatorID uint64, sig []byte) []byte {
	return append(append(bytesutil.Bytes8(uint64(slot)), bytesutil.Bytes8(validatorID)...), sig...)
}

func encodeEpochSig(targetEpoch basetypes.Epoch, sig []byte) []byte {
	return append(bytesutil.Bytes8(uint64(targetEpoch)), sig...)
}
func encodeType(st types.SlashingType) []byte {
	return []byte{byte(st)}
}
func encodeTypeRoot(st types.SlashingType, root [32]byte) []byte {
	return append([]byte{byte(st)}, root[:]...)
}
