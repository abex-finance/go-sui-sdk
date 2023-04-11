package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	SuiCoinType   = "0x2::sui::SUI"
	DevNetRpcUrl  = "https://fullnode.devnet.sui.io"
	TestnetRpcUrl = "https://fullnode.testnet.sui.io"
)

type Address = HexData

// NewAddressFromHex
/**
 * Creates Address from a hex string.
 * @param addr Hex string can be with a prefix or without a prefix,
 * e.g. '0x1aa' or '1aa'. Hex string will be left padded with 0s if too short.
 */
func NewAddressFromHex(addr string) (*Address, error) {
	if strings.HasPrefix(addr, "0x") || strings.HasPrefix(addr, "0X") {
		addr = addr[2:]
	}
	if len(addr)%2 != 0 {
		addr = "0" + addr
	}

	data, err := hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	const addressLength = 32
	if len(data) > addressLength {
		return nil, fmt.Errorf("hex string is too long. Address's length is %v data", addressLength)
	}

	res := [addressLength]byte{}
	copy(res[addressLength-len(data):], data[:])
	address := Address(res[:])
	return &address, nil
}

// ShortString Returns the address with leading zeros trimmed, e.g. 0x2
func (a Address) ShortString() string {
	return "0x" + strings.TrimLeft(hex.EncodeToString(a), "0")
}

type ObjectId = HexData
type Digest = Base64Data

type InputObjectKind map[string]interface{}

type TransactionBytes struct {
	// the gas object to be used
	Gas []ObjectRef `json:"gas"`

	// objects to be used in this transaction
	InputObjects []InputObjectKind `json:"inputObjects"`

	// transaction data bytes
	TxBytes Base64Data `json:"txBytes"`
}

// ObjectRef for BCS, need to keep this order
type ObjectRef struct {
	ObjectId ObjectId          `json:"objectId"`
	Version  SuiBigInt         `json:"version"`
	Digest   TransactionDigest `json:"digest"`
}

type TransferObject struct {
	Recipient Address   `json:"recipient"`
	ObjectRef ObjectRef `json:"object_ref"`
}
type ModulePublish struct {
	Modules [][]byte `json:"modules"`
}
type MoveCall struct {
	Package  ObjectId      `json:"package"`
	Module   string        `json:"module"`
	Function string        `json:"function"`
	TypeArgs []interface{} `json:"typeArguments"`
	Args     []interface{} `json:"arguments"`
}
type TransferSui struct {
	Recipient Address `json:"recipient"`
	Amount    uint64  `json:"amount"`
}
type Pay struct {
	Coins      []ObjectRef `json:"coins"`
	Recipients []Address   `json:"recipients"`
	Amounts    []uint64    `json:"amounts"`
}
type PaySui struct {
	Coins      []ObjectRef `json:"coins"`
	Recipients []Address   `json:"recipients"`
	Amounts    []uint64    `json:"amounts"`
}
type PayAllSui struct {
	Coins     []ObjectRef `json:"coins"`
	Recipient Address     `json:"recipient"`
}
type ChangeEpoch struct {
	Epoch             interface{} `json:"epoch"`
	StorageCharge     uint64      `json:"storage_charge"`
	ComputationCharge uint64      `json:"computation_charge"`
}

type SingleTransactionKind struct {
	TransferObject *TransferObject `json:"TransferObject,omitempty"`
	Publish        *ModulePublish  `json:"Publish,omitempty"`
	Call           *MoveCall       `json:"Call,omitempty"`
	TransferSui    *TransferSui    `json:"TransferSui,omitempty"`
	ChangeEpoch    *ChangeEpoch    `json:"ChangeEpoch,omitempty"`
	PaySui         *PaySui         `json:"PaySui,omitempty"`
	Pay            *Pay            `json:"Pay,omitempty"`
	PayAllSui      *PayAllSui      `json:"PayAllSui,omitempty"`
}

type SenderSignedData struct {
	Transactions []SingleTransactionKind `json:"transactions,omitempty"`

	Sender     *Address   `json:"sender"`
	GasPayment *ObjectRef `json:"gasPayment"`
	GasBudget  uint64     `json:"gasBudget"`
	// GasPrice     uint64      `json:"gasPrice"`
}

type TimeRange struct {
	StartTime uint64 `json:"startTime"` // left endpoint of time interval, milliseconds since epoch, inclusive
	EndTime   uint64 `json:"endTime"`   // right endpoint of time interval, milliseconds since epoch, exclusive
}

type MoveModule struct {
	Package ObjectId `json:"package"`
	Module  string   `json:"module"`
}

func (o ObjectOwner) MarshalJSON() ([]byte, error) {
	if o.string != nil {
		data, err := json.Marshal(o.string)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if o.ObjectOwnerInternal != nil {
		data, err := json.Marshal(o.ObjectOwnerInternal)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, errors.New("nil value")
}

func (o *ObjectOwner) UnmarshalJSON(data []byte) error {
	if bytes.HasPrefix(data, []byte("\"")) {
		stringData := string(data[1 : len(data)-1])
		o.string = &stringData
		return nil
	}
	if bytes.HasPrefix(data, []byte("{")) {
		oOI := ObjectOwnerInternal{}
		err := json.Unmarshal(data, &oOI)
		if err != nil {
			return err
		}
		o.ObjectOwnerInternal = &oOI
		return nil
	}
	return errors.New("value not json")
}

func IsSameStringAddress(addr1, addr2 string) bool {
	if strings.HasPrefix(addr1, "0x") {
		addr1 = addr1[2:]
	}
	if strings.HasPrefix(addr2, "0x") {
		addr2 = addr2[2:]
	}
	addr1 = strings.TrimLeft(addr1, "0")
	return strings.TrimLeft(addr1, "0") == strings.TrimLeft(addr2, "0")
}
