package infra

import (
	"encoding/binary"
	"protohackers/11_pest/types"
)

func Encode(v any) []byte {
	switch nv := v.(type) {
	case types.Hello:
		return encodeHello(nv)
	case types.Error:
		return encodeError(nv)
	case types.OK:
		return encodeOK(nv)
	case types.DialAuthority:
		return encodeDialAuthority(nv)
	case types.TargetPopulations:
		return encodeTargetPopulations(nv)
	case types.CreatePolicy:
		return encodeCreatePolicy(nv)
	case types.DeletePolicy:
		return encodeDeletePolicy(nv)
	case types.PolicyResult:
		return encodePolicyResult(nv)
	case types.SiteVisit:
		return encodeSiteVisit(nv)
	}
	return nil
}

func encodeHello(v types.Hello) []byte {
	body := make([]byte, 0)
	body = append(body, encodeString(v.Protocol)...)
	body = append(body, encodeUint32(v.Version)...)

	return encaseEnvelope(body, 0x50)
}

func encodeError(v types.Error) []byte {
	body := make([]byte, 0)
	body = append(body, encodeString(v.Message)...)
	return encaseEnvelope(body, 0x51)
}

func encodeOK(types.OK) []byte {
	return encaseEnvelope(nil, 0x52)
}

func encodeDialAuthority(v types.DialAuthority) []byte {
	body := make([]byte, 0)
	body = append(body, encodeUint32(v.Site)...)
	return encaseEnvelope(body, 0x53)
}

func encodeTargetPopulations(v types.TargetPopulations) []byte {
	body := make([]byte, 0)
	body = append(body, encodeUint32(v.Site)...)
	arrLen := uint32(len(v.Populations))
	body = append(body, encodeUint32(arrLen)...)
	for _, entry := range v.Populations {
		body = append(body, encodeString(entry.Species)...)
		body = append(body, encodeUint32(entry.Min)...)
		body = append(body, encodeUint32(entry.Max)...)
	}
	return encaseEnvelope(body, 0x54)
}

func encodeCreatePolicy(v types.CreatePolicy) []byte {
	body := make([]byte, 0)
	body = append(body, encodeString(v.Species)...)
	body = append(body, byte(v.Action))
	return encaseEnvelope(body, 0x55)
}

func encodeDeletePolicy(v types.DeletePolicy) []byte {
	body := make([]byte, 0)
	body = append(body, encodeUint32(v.Policy)...)
	return encaseEnvelope(body, 0x56)
}

func encodePolicyResult(v types.PolicyResult) []byte {
	body := make([]byte, 0)
	body = append(body, encodeUint32(v.Policy)...)
	return encaseEnvelope(body, 0x57)
}

func encodeSiteVisit(v types.SiteVisit) []byte {
	body := make([]byte, 0)
	body = append(body, encodeUint32(v.Site)...)
	arrLen := uint32(len(v.Populations))
	body = append(body, encodeUint32(arrLen)...)
	for _, entry := range v.Populations {
		body = append(body, encodeString(entry.Species)...)
		body = append(body, encodeUint32(entry.Count)...)
	}
	return encaseEnvelope(body, 0x58)
}

func encaseEnvelope(b []byte, prefix byte) (ret []byte) {
	// total len = prefix (1 byte) + msgLen (4 byte) + msg (len(b) bytes) + checksum (1 byte)
	// 1 + 4 + len(b) + 1
	// len(b) + 6
	totalLen := uint32(len(b) + 6)

	ret = append(ret, prefix)
	ret = append(ret, encodeUint32(totalLen)...)
	ret = append(ret, b...)

	var sum byte = 0
	for i := 0; i < len(ret); i++ {
		sum += ret[i]
	}
	checkSum := 255 - sum + 1
	ret = append(ret, checkSum)

	return ret
}

func encodeUint32(val uint32) (ret []byte) {
	ret = binary.BigEndian.AppendUint32(ret, val)
	return ret
}

func encodeString(s string) (ret []byte) {
	strlen := uint32(len(s))
	ret = append(ret, encodeUint32(strlen)...)
	ret = append(ret, []byte(s)...)
	return ret
}
