

// Code generated by "cofing -t BlockHeader"; DO NOT EDIT.
package message

import (
	"bytes"

	"github.com/t10471/bitcoin-coding/basetype"
)

func (b *BlockHeader) Decode(b_ *bytes.Buffer) error {
	// Version
	{
		var err error
		b.Version, err = basetype.DecodeInt32(b_)
		if err != nil {
			return err
		}
	} 
	// PrevBlock
	{
		var err error
		b.PrevBlock, err = basetype.DecodeHash(b_)
		if err != nil {
			return err
		}
	} 
	// MerkleRoot
	{
		var err error
		b.MerkleRoot, err = basetype.DecodeHash(b_)
		if err != nil {
			return err
		}
	} 
	// Timestamp
	{
		var err error
		b.Timestamp, err = basetype.DecodeUint32Time(b_)
		if err != nil {
			return err
		}
	} 
	// Bits
	{
		var err error
		b.Bits, err = basetype.DecodeUint32(b_)
		if err != nil {
			return err
		}
	} 
	// Nonce
	{
		var err error
		b.Nonce, err = basetype.DecodeUint32(b_)
		if err != nil {
			return err
		}
	}  
	return nil
}

func (b *BlockHeader) Encode(b_ *bytes.Buffer) error {
	// Version
	if err := basetype.EncodeInt32(b_, b.Version); err != nil {
		return err
	} 
	// PrevBlock
	if err := basetype.EncodeHash(b_, b.PrevBlock); err != nil {
		return err
	} 
	// MerkleRoot
	if err := basetype.EncodeHash(b_, b.MerkleRoot); err != nil {
		return err
	} 
	// Timestamp
	if err := basetype.EncodeUint32Time(b_, b.Timestamp); err != nil {
		return err
	} 
	// Bits
	if err := basetype.EncodeUint32(b_, b.Bits); err != nil {
		return err
	} 
	// Nonce
	if err := basetype.EncodeUint32(b_, b.Nonce); err != nil {
		return err
	}  
	return nil
}
