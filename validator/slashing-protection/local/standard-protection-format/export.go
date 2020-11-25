package interchangeformat

import (
	"context"
	"fmt"

	"github.com/prysmaticlabs/prysm/validator/db"
)

// ExportStandardProtectionJSON --
func ExportStandardProtectionJSON(ctx context.Context, validatorDB db.Database) (*EIPSlashingProtectionFormat, error) {
	interchangeJSON := &EIPSlashingProtectionFormat{}
	genesisValidatorsRoot, err := validatorDB.GenesisValidatorsRoot(ctx)
	if err != nil {
		return nil, err
	}
	genesisRootHex, err := rootToHex(genesisValidatorsRoot)
	if err != nil {
		return nil, err
	}
	interchangeJSON.Metadata.GenesisValidatorsRoot = genesisRootHex
	interchangeJSON.Metadata.InterchangeFormatVersion = INTERCHANGE_FORMAT_VERSION
	proposedPublicKeys, err := validatorDB.ProposedPublicKeys(ctx)
	if err != nil {
		return nil, err
	}
	data := make([]*ProtectionData, 0)
	for _, pubKey := range proposedPublicKeys {
		pubKeyHex, err := pubKeyToHex(pubKey[:])
		if err != nil {
			return nil, err
		}
		signedBlocks, err := getSignedBlocksByPubKey(ctx, validatorDB, pubKey)
		data = append(data, &ProtectionData{
			Pubkey:             pubKeyHex,
			SignedBlocks:       signedBlocks,
			SignedAttestations: nil,
		})
	}
	interchangeJSON.Data = data
	return interchangeJSON, nil
}

func getSignedBlocksByPubKey(ctx context.Context, validatorDB db.Database, pubKey [48]byte) ([]*SignedBlock, error) {
	lowestSignedSlot, err := validatorDB.LowestSignedProposal(ctx, pubKey[:])
	if err != nil {
		return nil, err
	}
	highestSignedSlot, err := validatorDB.HighestSignedProposal(ctx, pubKey[:])
	if err != nil {
		return nil, err
	}
	signedBlocks := make([]*SignedBlock, 0)
	for i := lowestSignedSlot; i < highestSignedSlot; i++ {
		signingRoot, err := validatorDB.ProposalHistoryForSlot(ctx, pubKey[:], i)
		if err != nil {
			return nil, err
		}
		signingRootHex, err := rootToHex(signingRoot)
		if err != nil {
			return nil, err
		}
		signedBlocks = append(signedBlocks, &SignedBlock{
			Slot:        fmt.Sprintf("%d", i),
			SigningRoot: signingRootHex,
		})
	}
	return signedBlocks, nil
}

func rootToHex(root []byte) (string, error) {
	if len(root) != 32 {
		return "", fmt.Errorf("wanted length 32, received %d", len(root))
	}
	return fmt.Sprintf("%#x", root), nil
}

func pubKeyToHex(pubKey []byte) (string, error) {
	if len(pubKey) != 48 {
		return "", fmt.Errorf("wanted length 48, received %d", len(pubKey))
	}
	return fmt.Sprintf("%#x", pubKey), nil
}