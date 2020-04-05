package basetype

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"time"
)

var (
	littleEndian = binary.LittleEndian
)

type (
	Hash       [32]byte
	VarInt     int
	Int32      int32
	Uint32     uint32
	Uint32Time time.Time
)

func DecodeUint32(buffer *bytes.Buffer) (Uint32, error) {
	v, err := decodeUint32(buffer, littleEndian)
	if err != nil {
		return 0, err
	}
	return Uint32(v), nil
}

func DecodeInt32(buffer *bytes.Buffer) (Int32, error) {
	rv, err := DecodeUint32(buffer)
	if err != nil {
		return 0, err
	}
	return Int32(int32(rv)), nil
}

func DecodeVarInt(buffer *bytes.Buffer) (VarInt, error) {
	v, err := decodeVarInt(buffer)
	if err != nil {
		return 0, err
	}
	return VarInt(v), nil
}

func DecodeHash(buffer *bytes.Buffer) (Hash, error) {
	var e Hash
	_, err := io.ReadFull(buffer, e[:])
	if err != nil {
		return e, err
	}
	return e, nil
}

func DecodeUint32Time(buffer *bytes.Buffer) (Uint32Time, error) {
	rv, err := DecodeUint32(buffer)
	if err != nil {
		return Uint32Time(time.Time{}), err
	}
	return Uint32Time(time.Unix(int64(rv), 0)), nil
}

func EncodeUint32(buffer *bytes.Buffer, val Uint32) error {
	return encodeUint32(buffer, littleEndian, uint32(val))
}
func EncodeInt32(buffer *bytes.Buffer, val Int32) error {
	return encodeUint32(buffer, littleEndian, uint32(val))
}

func EncodeVarInt(buffer *bytes.Buffer, val VarInt) error {
	return encodeVarInt(buffer, uint64(val))
}

func EncodeHash(buffer *bytes.Buffer, val Hash) error {
	_, err := buffer.Write(val[:])
	if err != nil {
		return err
	}
	return nil
}

func EncodeUint32Time(buffer *bytes.Buffer, val Uint32Time) error {
	return encodeUint32(buffer, littleEndian, uint32(time.Time(val).Unix()))
}

func DecodeVarInt(buffer *bytes.Buffer) (VarInt, error) {
	discriminant, err := decodeUint8(buffer)
	if err != nil {
		return 0, err
	}

	var rv uint64
	switch discriminant {
	case 0xff:
		rv, err := decodeUint64(buffer, littleEndian)
		if err != nil {
			return 0, err
		}
		if rv < uint64(0x100000000) {
			return 0, fmt.Errorf("non-canonical varint %x - discriminant %x must encode a value greater than %x", rv, discriminant, min)
		}

	case 0xfe:
		sv, err := decodeUint32(buffer, littleEndian)
		if err != nil {
			return 0, err
		}
		rv = uint64(sv)
		if rv < uint64(0x10000) {
			return 0, fmt.Errorf("non-canonical varint %x - discriminant %x must encode a value greater than %x", rv, discriminant, min)
		}

	case 0xfd:
		sv, err := decodeUint16(buffer, littleEndian)
		if err != nil {
			return 0, err
		}
		rv = uint64(sv)

		min := uint64(0xfd)
		if rv < min {
			return 0, fmt.Errorf("non-canonical varint %x - discriminant %x must encode a value greater than %x", rv, discriminant, min)
		}

	default:
		rv = uint64(discriminant)
	}
	return VarInt(rv), nil
}

func decodeVarInt(r io.Reader) (uint64, error) {
	discriminant, err := decodeUint8(r)
	if err != nil {
		return 0, err
	}

	var rv uint64
	switch discriminant {
	case 0xff:
		sv, err := decodeUint64(r, littleEndian)
		if err != nil {
			return 0, err
		}
		rv = sv

		min := uint64(0x100000000)
		if rv < min {
			return 0, fmt.Errorf("non-canonical varint %x - discriminant %x must encode a value greater than %x", rv, discriminant, min)
		}

	case 0xfe:
		sv, err := decodeUint32(r, littleEndian)
		if err != nil {
			return 0, err
		}
		rv = uint64(sv)

		min := uint64(0x10000)
		if rv < min {
			return 0, fmt.Errorf("non-canonical varint %x - discriminant %x must encode a value greater than %x", rv, discriminant, min)
		}

	case 0xfd:
		sv, err := decodeUint16(r, littleEndian)
		if err != nil {
			return 0, err
		}
		rv = uint64(sv)

		min := uint64(0xfd)
		if rv < min {
			return 0, fmt.Errorf("non-canonical varint %x - discriminant %x must encode a value greater than %x", rv, discriminant, min)
		}

	default:
		rv = uint64(discriminant)
	}

	return rv, nil
}

func decodeUint8(r io.Reader) (uint8, error) {
	buf := make([]byte, 1)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	rv := buf[0]
	return rv, nil
}

func decodeUint16(r io.Reader, byteOrder binary.ByteOrder) (uint16, error) {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	rv := byteOrder.Uint16(buf)
	return rv, nil
}

func decodeUint32(r io.Reader, byteOrder binary.ByteOrder) (uint32, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	rv := byteOrder.Uint32(buf)
	return rv, nil
}

func decodeUint64(r io.Reader, byteOrder binary.ByteOrder) (uint64, error) {
	buf := make([]byte, 8)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	rv := byteOrder.Uint64(buf)
	return rv, nil
}

func encodeVarInt(w io.Writer, val uint64) error {
	if val < 0xfd {
		return encodeUint8(w, uint8(val))
	}

	if val <= math.MaxUint16 {
		if err := encodeUint8(w, 0xfd); err != nil {
			return err
		}
		return encodeUint16(w, littleEndian, uint16(val))
	}

	if val <= math.MaxUint32 {
		if err := encodeUint8(w, 0xfe); err != nil {
			return err
		}
		return encodeUint32(w, littleEndian, uint32(val))
	}

	if err := encodeUint8(w, 0xff); err != nil {
		return err
	}

	return encodeUint64(w, littleEndian, val)
}

func encodeUint8(w io.Writer, val uint8) error {
	buf := make([]byte, 1)
	buf[0] = val
	_, err := w.Write(buf)
	return err
}

func encodeUint16(w io.Writer, byteOrder binary.ByteOrder, val uint16) error {
	buf := make([]byte, 2)
	byteOrder.PutUint16(buf, val)
	_, err := w.Write(buf)
	return err
}

func encodeUint32(w io.Writer, byteOrder binary.ByteOrder, val uint32) error {
	buf := make([]byte, 4)
	byteOrder.PutUint32(buf, val)
	_, err := w.Write(buf)
	return err
}

func encodeUint64(w io.Writer, byteOrder binary.ByteOrder, val uint64) error {
	buf := make([]byte, 8)
	byteOrder.PutUint64(buf, val)
	_, err := w.Write(buf)
	return err
}
