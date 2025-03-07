package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/flashbots/go-boost-utils/types"
)

func NewNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func NewNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

type GetPayloadsFilters struct {
	Slot           uint64
	Cursor         uint64
	Limit          uint64
	BlockHash      string
	BlockNumber    uint64
	ProposerPubkey string
	BuilderPubkey  string
	OrderByValue   int8
}

type GetBuilderSubmissionsFilters struct {
	Slot        uint64
	Limit       uint64
	BlockHash   string
	BlockNumber uint64
	// Cursor      uint64
	BuilderPubkey string
}

type ValidatorRegistrationEntry struct {
	ID         int64     `db:"id"`
	InsertedAt time.Time `db:"inserted_at"`

	Pubkey       string `db:"pubkey"`
	FeeRecipient string `db:"fee_recipient"`
	Timestamp    uint64 `db:"timestamp"`
	GasLimit     uint64 `db:"gas_limit"`
	Signature    string `db:"signature"`
}

func (reg ValidatorRegistrationEntry) ToSignedValidatorRegistration() (*types.SignedValidatorRegistration, error) {
	pubkey, err := types.HexToPubkey(reg.Pubkey)
	if err != nil {
		return nil, err
	}

	feeRec, err := types.HexToAddress(reg.FeeRecipient)
	if err != nil {
		return nil, err
	}

	sig, err := types.HexToSignature(reg.Signature)
	if err != nil {
		return nil, err
	}

	return &types.SignedValidatorRegistration{
		Message: &types.RegisterValidatorRequestMessage{
			Pubkey:       pubkey,
			FeeRecipient: feeRec,
			Timestamp:    reg.Timestamp,
			GasLimit:     reg.GasLimit,
		},
		Signature: sig,
	}, nil
}

func SignedValidatorRegistrationToEntry(valReg types.SignedValidatorRegistration) ValidatorRegistrationEntry {
	return ValidatorRegistrationEntry{
		Pubkey:       valReg.Message.Pubkey.String(),
		FeeRecipient: valReg.Message.FeeRecipient.String(),
		Timestamp:    valReg.Message.Timestamp,
		GasLimit:     valReg.Message.GasLimit,
		Signature:    valReg.Signature.String(),
	}
}

type ExecutionPayloadEntry struct {
	ID         int64     `db:"id"`
	InsertedAt time.Time `db:"inserted_at"`

	Slot           uint64 `db:"slot"`
	ProposerPubkey string `db:"proposer_pubkey"`
	BlockHash      string `db:"block_hash"`

	Version string `db:"version"`
	Payload string `db:"payload"`
}

var ExecutionPayloadEntryCSVHeader = []string{"id", "inserted_at", "slot", "proposer_pubkey", "block_hash", "version", "payload"}

func (e *ExecutionPayloadEntry) ToCSVRecord() []string {
	return []string{
		fmt.Sprint(e.ID),
		e.InsertedAt.UTC().String(),
		fmt.Sprint(e.Slot),
		e.ProposerPubkey,
		e.BlockHash,
		e.Version,
		e.Payload,
	}
}

type BuilderBlockSubmissionEntry struct {
	ID         int64        `db:"id"`
	InsertedAt time.Time    `db:"inserted_at"`
	ReceivedAt sql.NullTime `db:"received_at"`

	// Delivered ExecutionPayload
	ExecutionPayloadID sql.NullInt64 `db:"execution_payload_id"`

	// Sim Result
	SimSuccess bool   `db:"sim_success"`
	SimError   string `db:"sim_error"`

	// BidTrace data
	Signature string `db:"signature"`

	Slot       uint64 `db:"slot"`
	ParentHash string `db:"parent_hash"`
	BlockHash  string `db:"block_hash"`

	BuilderPubkey        string `db:"builder_pubkey"`
	ProposerPubkey       string `db:"proposer_pubkey"`
	ProposerFeeRecipient string `db:"proposer_fee_recipient"`

	GasUsed  uint64 `db:"gas_used"`
	GasLimit uint64 `db:"gas_limit"`

	NumTx uint64 `db:"num_tx"`
	Value string `db:"value"`

	// Helpers
	Epoch       uint64 `db:"epoch"`
	BlockNumber uint64 `db:"block_number"`
}

type DeliveredPayloadEntry struct {
	ID         int64     `db:"id"`
	InsertedAt time.Time `db:"inserted_at"`

	SignedBlindedBeaconBlock sql.NullString `db:"signed_blinded_beacon_block"`

	Slot  uint64 `db:"slot"`
	Epoch uint64 `db:"epoch"`

	BuilderPubkey        string `db:"builder_pubkey"`
	ProposerPubkey       string `db:"proposer_pubkey"`
	ProposerFeeRecipient string `db:"proposer_fee_recipient"`

	ParentHash  string `db:"parent_hash"`
	BlockHash   string `db:"block_hash"`
	BlockNumber uint64 `db:"block_number"`

	GasUsed  uint64 `db:"gas_used"`
	GasLimit uint64 `db:"gas_limit"`

	NumTx uint64 `db:"num_tx"`
	Value string `db:"value"`
}

type BlockBuilderEntry struct {
	ID         int64     `db:"id"          json:"id"`
	InsertedAt time.Time `db:"inserted_at" json:"inserted_at"`

	BuilderPubkey string `db:"builder_pubkey" json:"builder_pubkey"`
	Description   string `db:"description"    json:"description"`

	IsHighPrio    bool `db:"is_high_prio"   json:"is_high_prio"`
	IsBlacklisted bool `db:"is_blacklisted" json:"is_blacklisted"`

	LastSubmissionID   sql.NullInt64 `db:"last_submission_id"   json:"last_submission_id"`
	LastSubmissionSlot uint64        `db:"last_submission_slot" json:"last_submission_slot"`

	NumSubmissionsTotal    uint64 `db:"num_submissions_total"    json:"num_submissions_total"`
	NumSubmissionsSimError uint64 `db:"num_submissions_simerror" json:"num_submissions_simerror"`

	NumSentGetPayload uint64 `db:"num_sent_getpayload" json:"num_sent_getpayload"`
}
