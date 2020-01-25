// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ring.proto

package ring

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_sortkeys "github.com/gogo/protobuf/sortkeys"
	io "io"
	math "math"
	reflect "reflect"
	strconv "strconv"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type IngesterState int32

const (
	ACTIVE  IngesterState = 0
	LEAVING IngesterState = 1
	PENDING IngesterState = 2
	JOINING IngesterState = 3
	// This state is only used by gossiping code to distribute information about
	// ingesters that have been removed from the ring. Ring users should not use it directly.
	LEFT IngesterState = 4
)

var IngesterState_name = map[int32]string{
	0: "ACTIVE",
	1: "LEAVING",
	2: "PENDING",
	3: "JOINING",
	4: "LEFT",
}

var IngesterState_value = map[string]int32{
	"ACTIVE":  0,
	"LEAVING": 1,
	"PENDING": 2,
	"JOINING": 3,
	"LEFT":    4,
}

func (IngesterState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_26381ed67e202a6e, []int{0}
}

type Desc struct {
	Ingesters map[string]IngesterDesc `protobuf:"bytes,1,rep,name=ingesters,proto3" json:"ingesters" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *Desc) Reset()      { *m = Desc{} }
func (*Desc) ProtoMessage() {}
func (*Desc) Descriptor() ([]byte, []int) {
	return fileDescriptor_26381ed67e202a6e, []int{0}
}
func (m *Desc) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Desc) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Desc.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Desc) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Desc.Merge(m, src)
}
func (m *Desc) XXX_Size() int {
	return m.Size()
}
func (m *Desc) XXX_DiscardUnknown() {
	xxx_messageInfo_Desc.DiscardUnknown(m)
}

var xxx_messageInfo_Desc proto.InternalMessageInfo

func (m *Desc) GetIngesters() map[string]IngesterDesc {
	if m != nil {
		return m.Ingesters
	}
	return nil
}

type IngesterDesc struct {
	Addr      string        `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Timestamp int64         `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	State     IngesterState `protobuf:"varint,3,opt,name=state,proto3,enum=ring.IngesterState" json:"state,omitempty"`
	Tokens    []uint32      `protobuf:"varint,6,rep,packed,name=tokens,proto3" json:"tokens,omitempty"`
}

func (m *IngesterDesc) Reset()      { *m = IngesterDesc{} }
func (*IngesterDesc) ProtoMessage() {}
func (*IngesterDesc) Descriptor() ([]byte, []int) {
	return fileDescriptor_26381ed67e202a6e, []int{1}
}
func (m *IngesterDesc) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IngesterDesc) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IngesterDesc.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IngesterDesc) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IngesterDesc.Merge(m, src)
}
func (m *IngesterDesc) XXX_Size() int {
	return m.Size()
}
func (m *IngesterDesc) XXX_DiscardUnknown() {
	xxx_messageInfo_IngesterDesc.DiscardUnknown(m)
}

var xxx_messageInfo_IngesterDesc proto.InternalMessageInfo

func (m *IngesterDesc) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *IngesterDesc) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *IngesterDesc) GetState() IngesterState {
	if m != nil {
		return m.State
	}
	return ACTIVE
}

func (m *IngesterDesc) GetTokens() []uint32 {
	if m != nil {
		return m.Tokens
	}
	return nil
}

func init() {
	proto.RegisterEnum("ring.IngesterState", IngesterState_name, IngesterState_value)
	proto.RegisterType((*Desc)(nil), "ring.Desc")
	proto.RegisterMapType((map[string]IngesterDesc)(nil), "ring.Desc.IngestersEntry")
	proto.RegisterType((*IngesterDesc)(nil), "ring.IngesterDesc")
}

func init() { proto.RegisterFile("ring.proto", fileDescriptor_26381ed67e202a6e) }

var fileDescriptor_26381ed67e202a6e = []byte{
	// 382 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x91, 0x4f, 0x8b, 0xd3, 0x40,
	0x18, 0xc6, 0xe7, 0x6d, 0x26, 0xb1, 0x7d, 0xeb, 0x2e, 0x61, 0x04, 0x89, 0x8b, 0x8c, 0x61, 0x4f,
	0x51, 0x30, 0x0b, 0xd5, 0x83, 0x08, 0x1e, 0x76, 0xdd, 0x28, 0x29, 0xa5, 0x2e, 0x71, 0xd9, 0x7b,
	0xda, 0x8e, 0x31, 0xd4, 0x26, 0x25, 0x99, 0x0a, 0xbd, 0xf9, 0x0d, 0xf4, 0xe6, 0x57, 0xf0, 0xa3,
	0xf4, 0xd8, 0x63, 0x4f, 0x62, 0xd3, 0x8b, 0xc7, 0x7e, 0x04, 0x99, 0x49, 0x4b, 0xed, 0xed, 0xf9,
	0xcd, 0xf3, 0xe7, 0x3d, 0x0c, 0x62, 0x91, 0x66, 0x89, 0x3f, 0x2d, 0x72, 0x99, 0x33, 0xaa, 0xf4,
	0xd9, 0xf3, 0x24, 0x95, 0x9f, 0x67, 0x03, 0x7f, 0x98, 0x4f, 0x2e, 0x92, 0x3c, 0xc9, 0x2f, 0xb4,
	0x39, 0x98, 0x7d, 0xd2, 0xa4, 0x41, 0xab, 0xba, 0x74, 0xfe, 0x13, 0x90, 0x5e, 0x8b, 0x72, 0xc8,
	0xde, 0x60, 0x2b, 0xcd, 0x12, 0x51, 0x4a, 0x51, 0x94, 0x0e, 0xb8, 0x86, 0xd7, 0xee, 0x3c, 0xf2,
	0xf5, 0xba, 0xb2, 0xfd, 0x70, 0xef, 0x05, 0x99, 0x2c, 0xe6, 0x57, 0x74, 0xf1, 0xfb, 0x09, 0x89,
	0x0e, 0x8d, 0xb3, 0x1b, 0x3c, 0x3d, 0x8e, 0x30, 0x1b, 0x8d, 0xb1, 0x98, 0x3b, 0xe0, 0x82, 0xd7,
	0x8a, 0x94, 0x64, 0x1e, 0x9a, 0x5f, 0xe3, 0x2f, 0x33, 0xe1, 0x34, 0x5c, 0xf0, 0xda, 0x1d, 0x56,
	0xcf, 0xef, 0x6b, 0xea, 0x4c, 0x54, 0x07, 0x5e, 0x37, 0x5e, 0xc1, 0xf9, 0x77, 0xc0, 0xfb, 0xff,
	0x7b, 0x8c, 0x21, 0x8d, 0x47, 0xa3, 0x62, 0xb7, 0xa8, 0x35, 0x7b, 0x8c, 0x2d, 0x99, 0x4e, 0x44,
	0x29, 0xe3, 0xc9, 0x54, 0xcf, 0x1a, 0xd1, 0xe1, 0x81, 0x3d, 0x45, 0xb3, 0x94, 0xb1, 0x14, 0x8e,
	0xe1, 0x82, 0x77, 0xda, 0x79, 0x70, 0x7c, 0xf0, 0xa3, 0xb2, 0xa2, 0x3a, 0xc1, 0x1e, 0xa2, 0x25,
	0xf3, 0xb1, 0xc8, 0x4a, 0xc7, 0x72, 0x0d, 0xef, 0x24, 0xda, 0x51, 0x97, 0x36, 0xa9, 0x6d, 0x76,
	0x69, 0xd3, 0xb4, 0xad, 0x67, 0x3d, 0x3c, 0x39, 0xea, 0x32, 0x44, 0xeb, 0xf2, 0xed, 0x6d, 0x78,
	0x17, 0xd8, 0x84, 0xb5, 0xf1, 0x5e, 0x2f, 0xb8, 0xbc, 0x0b, 0xfb, 0xef, 0x6d, 0x50, 0x70, 0x13,
	0xf4, 0xaf, 0x15, 0x34, 0x14, 0x74, 0x3f, 0x84, 0x7d, 0x05, 0x06, 0x6b, 0x22, 0xed, 0x05, 0xef,
	0x6e, 0x6d, 0x7a, 0xf5, 0x72, 0xb9, 0xe6, 0x64, 0xb5, 0xe6, 0x64, 0xbb, 0xe6, 0xf0, 0xad, 0xe2,
	0xf0, 0xab, 0xe2, 0xb0, 0xa8, 0x38, 0x2c, 0x2b, 0x0e, 0x7f, 0x2a, 0x0e, 0x7f, 0x2b, 0x4e, 0xb6,
	0x15, 0x87, 0x1f, 0x1b, 0x4e, 0x96, 0x1b, 0x4e, 0x56, 0x1b, 0x4e, 0x06, 0x96, 0xfe, 0xb6, 0x17,
	0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x39, 0x00, 0x90, 0xda, 0xf9, 0x01, 0x00, 0x00,
}

func (x IngesterState) String() string {
	s, ok := IngesterState_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (this *Desc) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Desc)
	if !ok {
		that2, ok := that.(Desc)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.Ingesters) != len(that1.Ingesters) {
		return false
	}
	for i := range this.Ingesters {
		a := this.Ingesters[i]
		b := that1.Ingesters[i]
		if !(&a).Equal(&b) {
			return false
		}
	}
	return true
}
func (this *IngesterDesc) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*IngesterDesc)
	if !ok {
		that2, ok := that.(IngesterDesc)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Addr != that1.Addr {
		return false
	}
	if this.Timestamp != that1.Timestamp {
		return false
	}
	if this.State != that1.State {
		return false
	}
	if len(this.Tokens) != len(that1.Tokens) {
		return false
	}
	for i := range this.Tokens {
		if this.Tokens[i] != that1.Tokens[i] {
			return false
		}
	}
	return true
}
func (this *Desc) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&ring.Desc{")
	keysForIngesters := make([]string, 0, len(this.Ingesters))
	for k, _ := range this.Ingesters {
		keysForIngesters = append(keysForIngesters, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForIngesters)
	mapStringForIngesters := "map[string]IngesterDesc{"
	for _, k := range keysForIngesters {
		mapStringForIngesters += fmt.Sprintf("%#v: %#v,", k, this.Ingesters[k])
	}
	mapStringForIngesters += "}"
	if this.Ingesters != nil {
		s = append(s, "Ingesters: "+mapStringForIngesters+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *IngesterDesc) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&ring.IngesterDesc{")
	s = append(s, "Addr: "+fmt.Sprintf("%#v", this.Addr)+",\n")
	s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	s = append(s, "State: "+fmt.Sprintf("%#v", this.State)+",\n")
	s = append(s, "Tokens: "+fmt.Sprintf("%#v", this.Tokens)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringRing(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *Desc) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Desc) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Ingesters) > 0 {
		for k, _ := range m.Ingesters {
			dAtA[i] = 0xa
			i++
			v := m.Ingesters[k]
			msgSize := 0
			if (&v) != nil {
				msgSize = (&v).Size()
				msgSize += 1 + sovRing(uint64(msgSize))
			}
			mapSize := 1 + len(k) + sovRing(uint64(len(k))) + msgSize
			i = encodeVarintRing(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintRing(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintRing(dAtA, i, uint64((&v).Size()))
			n1, err := (&v).MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n1
		}
	}
	return i, nil
}

func (m *IngesterDesc) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IngesterDesc) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Addr) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintRing(dAtA, i, uint64(len(m.Addr)))
		i += copy(dAtA[i:], m.Addr)
	}
	if m.Timestamp != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintRing(dAtA, i, uint64(m.Timestamp))
	}
	if m.State != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintRing(dAtA, i, uint64(m.State))
	}
	if len(m.Tokens) > 0 {
		dAtA3 := make([]byte, len(m.Tokens)*10)
		var j2 int
		for _, num := range m.Tokens {
			for num >= 1<<7 {
				dAtA3[j2] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j2++
			}
			dAtA3[j2] = uint8(num)
			j2++
		}
		dAtA[i] = 0x32
		i++
		i = encodeVarintRing(dAtA, i, uint64(j2))
		i += copy(dAtA[i:], dAtA3[:j2])
	}
	return i, nil
}

func encodeVarintRing(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Desc) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Ingesters) > 0 {
		for k, v := range m.Ingesters {
			_ = k
			_ = v
			l = v.Size()
			mapEntrySize := 1 + len(k) + sovRing(uint64(len(k))) + 1 + l + sovRing(uint64(l))
			n += mapEntrySize + 1 + sovRing(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *IngesterDesc) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Addr)
	if l > 0 {
		n += 1 + l + sovRing(uint64(l))
	}
	if m.Timestamp != 0 {
		n += 1 + sovRing(uint64(m.Timestamp))
	}
	if m.State != 0 {
		n += 1 + sovRing(uint64(m.State))
	}
	if len(m.Tokens) > 0 {
		l = 0
		for _, e := range m.Tokens {
			l += sovRing(uint64(e))
		}
		n += 1 + sovRing(uint64(l)) + l
	}
	return n
}

func sovRing(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozRing(x uint64) (n int) {
	return sovRing(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Desc) String() string {
	if this == nil {
		return "nil"
	}
	keysForIngesters := make([]string, 0, len(this.Ingesters))
	for k, _ := range this.Ingesters {
		keysForIngesters = append(keysForIngesters, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForIngesters)
	mapStringForIngesters := "map[string]IngesterDesc{"
	for _, k := range keysForIngesters {
		mapStringForIngesters += fmt.Sprintf("%v: %v,", k, this.Ingesters[k])
	}
	mapStringForIngesters += "}"
	s := strings.Join([]string{`&Desc{`,
		`Ingesters:` + mapStringForIngesters + `,`,
		`}`,
	}, "")
	return s
}
func (this *IngesterDesc) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&IngesterDesc{`,
		`Addr:` + fmt.Sprintf("%v", this.Addr) + `,`,
		`Timestamp:` + fmt.Sprintf("%v", this.Timestamp) + `,`,
		`State:` + fmt.Sprintf("%v", this.State) + `,`,
		`Tokens:` + fmt.Sprintf("%v", this.Tokens) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringRing(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Desc) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRing
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Desc: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Desc: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ingesters", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRing
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRing
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRing
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Ingesters == nil {
				m.Ingesters = make(map[string]IngesterDesc)
			}
			var mapkey string
			mapvalue := &IngesterDesc{}
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowRing
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowRing
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthRing
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthRing
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowRing
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthRing
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthRing
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &IngesterDesc{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipRing(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthRing
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Ingesters[mapkey] = *mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRing(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRing
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRing
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *IngesterDesc) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRing
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: IngesterDesc: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IngesterDesc: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Addr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRing
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRing
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRing
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Addr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRing
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRing
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= IngesterState(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType == 0 {
				var v uint32
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowRing
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint32(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Tokens = append(m.Tokens, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowRing
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthRing
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthRing
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.Tokens) == 0 {
					m.Tokens = make([]uint32, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint32
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowRing
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint32(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Tokens = append(m.Tokens, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Tokens", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRing(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRing
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRing
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipRing(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRing
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRing
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRing
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthRing
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthRing
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowRing
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipRing(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthRing
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthRing = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRing   = fmt.Errorf("proto: integer overflow")
)
