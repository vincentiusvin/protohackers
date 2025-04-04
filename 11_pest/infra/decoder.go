package infra

import (
	"encoding/binary"
	"fmt"
	"protohackers/11_pest/types"
)

var (
	ErrInvalidChecksum = fmt.Errorf("invalid checksum")
	ErrInvalidLength   = fmt.Errorf("invalid message length")
	ErrInvalidPrefix   = fmt.Errorf("invalid prefix")
	ErrInvalidData     = fmt.Errorf("invalid data")
	ErrNotEnough       = fmt.Errorf("not enough data")
	ErrTooLong         = fmt.Errorf("length too long")
)

type ParseFunc[T any] func(b []byte) ParseResult[T]
type ParseResult[T any] struct {
	Value T
	Next  []byte
	Error error
}

func Parse(b []byte) (ret ParseResult[any]) {
	// parse prefix first, but do not advance the byte stream.
	prefix := parseUint8(b)
	if prefix.Error != nil {
		ret.Error = ErrNotEnough
		return ret
	}

	switch prefix.Value {
	case 0x50:
		res := parseHello(b)
		ret.Error = res.Error
		ret.Next = res.Next
		ret.Value = res.Value
	case 0x51:
		res := parseError(b)
		ret.Error = res.Error
		ret.Next = res.Next
		ret.Value = res.Value
	case 0x52:
		res := parseOk(b)
		ret.Error = res.Error
		ret.Next = res.Next
		ret.Value = res.Value
	case 0x53:
		res := parseDialAuthority(b)
		ret.Error = res.Error
		ret.Next = res.Next
		ret.Value = res.Value
	case 0x54:
		res := parseTargetPopulations(b)
		ret.Error = res.Error
		ret.Next = res.Next
		ret.Value = res.Value
	case 0x55:
		res := parseCreatePolicy(b)
		ret.Error = res.Error
		ret.Next = res.Next
		ret.Value = res.Value
	case 0x56:
		res := parseDeletePolicy(b)
		ret.Error = res.Error
		ret.Next = res.Next
		ret.Value = res.Value
	case 0x57:
		res := parsePolicyResult(b)
		ret.Error = res.Error
		ret.Next = res.Next
		ret.Value = res.Value
	case 0x58:
		res := parseSiteVisit(b)
		ret.Error = res.Error
		ret.Next = res.Next
		ret.Value = res.Value
	default:
		ret.Error = ErrInvalidPrefix
		ret.Value = prefix.Value
		ret.Next = prefix.Next
	}

	return ret
}

func parseHello(b []byte) ParseResult[types.Hello] {
	return envelope(func(b []byte) (ret ParseResult[types.Hello]) {
		protocol := parseStringLimit(b, 11)
		if protocol.Error != nil {
			ret.Error = protocol.Error
			return ret
		}
		version := parseUint32(protocol.Next)
		if version.Error != nil {
			ret.Error = version.Error
			return ret
		}

		ret.Value = types.Hello{
			Protocol: protocol.Value,
			Version:  version.Value,
		}
		ret.Next = version.Next
		return ret
	}, 0x50)(b)
}

func parseError(b []byte) ParseResult[types.Error] {
	return envelope(func(b []byte) (ret ParseResult[types.Error]) {
		message := parseString(b)
		if message.Error != nil {
			ret.Error = message.Error
			return ret
		}
		ret.Value = types.Error{
			Message: message.Value,
		}
		ret.Next = message.Next
		return ret
	}, 0x51)(b)
}

func parseOk(b []byte) ParseResult[types.OK] {
	return envelope(func(b []byte) (ret ParseResult[types.OK]) {
		ret.Next = b
		return ret
	}, 0x52)(b)
}

func parseDialAuthority(b []byte) ParseResult[types.DialAuthority] {
	return envelope(func(b []byte) (ret ParseResult[types.DialAuthority]) {
		site := parseUint32(b)
		if site.Error != nil {
			ret.Error = site.Error
			return ret
		}

		ret.Value = types.DialAuthority{
			Site: site.Value,
		}
		ret.Next = site.Next
		return ret
	}, 0x53)(b)
}

func parseTargetPopulationsEntry(b []byte) (ret ParseResult[types.TargetPopulationsEntry]) {
	species := parseString(b)
	if species.Error != nil {
		ret.Error = species.Error
		return ret
	}

	min := parseUint32(species.Next)
	if min.Error != nil {
		ret.Error = min.Error
		return ret
	}

	max := parseUint32(min.Next)
	if max.Error != nil {
		ret.Error = min.Error
		return ret
	}

	ret.Value = types.TargetPopulationsEntry{
		Species: species.Value,
		Min:     min.Value,
		Max:     max.Value,
	}
	ret.Next = max.Next
	return ret
}

func parseTargetPopulations(b []byte) ParseResult[types.TargetPopulations] {
	return envelope(func(b []byte) (ret ParseResult[types.TargetPopulations]) {
		site := parseUint32(b)
		if site.Error != nil {
			ret.Error = site.Error
			return ret
		}
		pops := parseArray(parseTargetPopulationsEntry)(site.Next)
		if pops.Error != nil {
			ret.Error = site.Error
			return ret
		}

		ret.Value = types.TargetPopulations{
			Site:        site.Value,
			Populations: pops.Value,
		}
		ret.Next = pops.Next
		return ret
	}, 0x54)(b)
}

func parseCreatePolicy(b []byte) ParseResult[types.CreatePolicy] {
	return envelope(func(b []byte) (ret ParseResult[types.CreatePolicy]) {
		species := parseString(b)
		if species.Error != nil {
			ret.Error = species.Error
			return ret
		}
		action := parseUint8(species.Next)
		if action.Error != nil {
			ret.Error = action.Error
			return ret
		}

		actionValue := types.Policy(action.Value)
		if actionValue != types.PolicyCull && actionValue != types.PolicyConserve {
			ret.Error = ErrInvalidData
			return ret
		}

		ret.Value = types.CreatePolicy{
			Species: species.Value,
			Action:  actionValue,
		}
		ret.Next = action.Next
		return ret
	}, 0x55)(b)
}

func parseDeletePolicy(b []byte) ParseResult[types.DeletePolicy] {
	return envelope(func(b []byte) (ret ParseResult[types.DeletePolicy]) {
		policy := parseUint32(b)
		if policy.Error != nil {
			ret.Error = policy.Error
			return ret
		}

		ret.Value = types.DeletePolicy{
			Policy: policy.Value,
		}
		ret.Next = policy.Next
		return ret
	}, 0x56)(b)
}

func parsePolicyResult(b []byte) ParseResult[types.PolicyResult] {
	return envelope(func(b []byte) (ret ParseResult[types.PolicyResult]) {
		policy := parseUint32(b)
		if policy.Error != nil {
			ret.Error = policy.Error
			return ret
		}

		ret.Value = types.PolicyResult{
			Policy: policy.Value,
		}
		ret.Next = policy.Next
		return ret
	}, 0x57)(b)
}

func parseSiteVisitEntry(b []byte) (ret ParseResult[types.SiteVisitEntry]) {
	species := parseString(b)
	if species.Error != nil {
		ret.Error = species.Error
		return ret
	}

	count := parseUint32(species.Next)
	if count.Error != nil {
		ret.Error = count.Error
		return ret
	}

	ret.Value = types.SiteVisitEntry{
		Species: species.Value,
		Count:   count.Value,
	}
	ret.Next = count.Next
	return ret
}

func parseSiteVisit(b []byte) ParseResult[types.SiteVisit] {
	msgLen := parseUint32(b[1:])
	if msgLen.Error != nil {
		return ParseResult[types.SiteVisit]{
			Error: msgLen.Error,
		}
	}

	return envelope(func(b []byte) (ret ParseResult[types.SiteVisit]) {
		site := parseUint32(b)
		if site.Error != nil {
			ret.Error = site.Error
			return ret
		}

		// total = prefix (1) + msgLen(4) + site(4) + arr(x) + checksum(1)
		// arr = total - 10
		limit := (int(msgLen.Value) - 10) / 1

		pops := parseArrayLimit(parseSiteVisitEntry, limit)(site.Next)
		if pops.Error != nil {
			ret.Error = pops.Error
			return ret
		}

		ret.Value = types.SiteVisit{
			Site:        site.Value,
			Populations: pops.Value,
		}
		ret.Next = pops.Next
		return ret
	}, 0x58)(b)
}

// envelopes the parser function fn with:
// - prefix verification
// - message length verification.
// - checksum verification
func envelope[T any](fn ParseFunc[T], expectedPrefix uint8) ParseFunc[T] {
	return func(b []byte) (ret ParseResult[T]) {
		prefix := parseUint8(b)
		if prefix.Error != nil {
			ret.Error = prefix.Error
			return ret
		}

		if prefix.Value != expectedPrefix {
			ret.Error = ErrInvalidData
			return ret
		}

		msgLen := parseUint32(prefix.Next)
		if msgLen.Error != nil {
			ret.Error = msgLen.Error
			return ret
		}
		val := fn(msgLen.Next)
		if val.Error != nil {
			ret.Error = val.Error
			return ret
		}

		checksum := parseUint8(val.Next)
		if checksum.Error != nil {
			ret.Error = checksum.Error
			return ret
		}

		expectedMsgLen := int(msgLen.Value)
		actualMsgLen := len(b) - len(checksum.Next)
		if expectedMsgLen != actualMsgLen {
			ret.Error = ErrInvalidLength
			return ret
		}

		var sum uint8
		for i := 0; i < actualMsgLen; i++ {
			sum += b[i]
		}
		if sum != 0 {
			ret.Error = ErrInvalidChecksum
			return ret
		}

		ret.Value = val.Value
		ret.Next = checksum.Next
		return ret
	}
}

// parser combinator :)
func parseArray[T any](fn ParseFunc[T]) ParseFunc[[]T] {
	return func(b []byte) (ret ParseResult[[]T]) {
		lenParse := parseUint32(b)
		if lenParse.Error != nil {
			ret.Error = lenParse.Error
			return ret
		}

		b = lenParse.Next
		lenVal := int(lenParse.Value)

		acc := make([]T, lenVal)
		for i := 0; i < lenVal; i++ {
			curr := fn(b)
			if curr.Error != nil {
				ret.Error = curr.Error
				return ret
			}
			acc[i] = curr.Value
			b = curr.Next
		}

		ret.Value = acc
		ret.Next = b

		return ret
	}
}

func parseArrayLimit[T any](fn ParseFunc[T], byteLimit int) ParseFunc[[]T] {
	return func(b []byte) (ret ParseResult[[]T]) {
		init := b
		lenParse := parseUint32(b)
		if lenParse.Error != nil {
			ret.Error = lenParse.Error
			return ret
		}

		b = lenParse.Next
		lenVal := int(lenParse.Value)

		acc := make([]T, lenVal)
		for i := 0; i < lenVal; i++ {
			curr := fn(b)
			if curr.Error != nil {
				ret.Error = curr.Error
				return ret
			}
			acc[i] = curr.Value
			b = curr.Next

			diff := len(init) - len(b)
			if diff >= byteLimit && (i+1 < lenVal) {
				ret.Error = ErrTooLong
				return
			}
		}

		ret.Value = acc
		ret.Next = b

		return ret
	}
}

func parseUint8(b []byte) (ret ParseResult[uint8]) {
	if len(b) < 1 {
		ret.Error = ErrNotEnough
		return
	}

	ret.Value = uint8(b[0])
	ret.Next = b[1:]
	return ret
}

func parseUint32(b []byte) (ret ParseResult[uint32]) {
	if len(b) < 4 {
		ret.Error = ErrNotEnough
		return
	}

	ret.Value = binary.BigEndian.Uint32(b)
	ret.Next = b[4:]
	return ret
}

// Consumes tokens from b to produce a string
// Returns number of bytes consumed and the final string
func parseString(b []byte) (ret ParseResult[string]) {
	lenParse := parseUint32(b)
	if lenParse.Error != nil {
		ret.Error = ErrNotEnough
		return
	}

	b = lenParse.Next
	lenVal := int(lenParse.Value)

	if len(b) < lenVal {
		ret.Error = ErrNotEnough
		return ret
	}

	str := string(b[:lenVal])

	ret.Value = str
	ret.Next = b[lenVal:]

	return ret
}

func parseStringLimit(b []byte, limit int) (ret ParseResult[string]) {
	lenParse := parseUint32(b)
	if lenParse.Error != nil {
		ret.Error = ErrNotEnough
		return
	}

	lenVal := int(lenParse.Value)
	if lenVal > limit {
		ret.Error = ErrTooLong
		return
	}

	return parseString(b)
}
