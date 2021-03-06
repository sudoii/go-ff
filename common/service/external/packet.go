package external

import (
	"encoding/binary"
	"math"
	"strings"

	"github.com/golang/geo/r3"
)

// Packet stores Data & Offset for a in/out message
type Packet struct {
	Data             []byte
	Offset           uint32
	MergePacketCount uint16
}

// MakePacket constructs a new packet with the given protocol
func MakePacket(protocol uint32) *Packet {
	packet := new(Packet)
	packet.Data = make([]byte, 10240)
	packet.MergePacketCount = 0

	return packet.
		WriteUInt8(0x5E).
		WriteUInt32(0).
		WriteUInt32(protocol)
}

// StartMergePacket create a simple packet without protocol & Data
func StartMergePacket(moverID uint32, cmd uint16, mainCmd uint32) *Packet {
	packet := new(Packet)
	packet.Data = make([]byte, 10240)
	packet.MergePacketCount = 1

	return packet.
		WriteUInt8(0x5E).
		WriteUInt32(0).
		WriteUInt32(mainCmd).
		WriteUInt32(0).
		WriteUInt16(packet.MergePacketCount). // MergePacketCount
		WriteUInt32(moverID).
		WriteUInt16(uint16(cmd))
}

// AddMergePart add a new entry to the packet
func (p *Packet) AddMergePart(protocol uint16, moverID uint32) *Packet {
	p.MergePacketCount++

	lastOffset := p.Offset
	p.Offset = 13
	p.WriteUInt16(p.MergePacketCount)
	p.Offset = lastOffset

	p.WriteUInt32(moverID).
		WriteUInt16(uint16(protocol))

	return p
}

// Finalize prepares the packet to be send by adding size
func (p *Packet) Finalize() *Packet {
	totalLen := p.Offset
	if totalLen < 5 {
		return p
	}

	binary.LittleEndian.PutUint32(p.Data[1:], uint32(totalLen-5))
	return p
}

// FinalizeForInternal prepares the packet to be sent internally by adding size
func (p *Packet) FinalizeForInternal() *Packet {
	totalLen := p.Offset
	if totalLen < 5 {
		return p
	}

	binary.LittleEndian.PutUint32(p.Data[1:], uint32(totalLen-5))
	p.Offset -= (32 / 8)
	return p
}

// WriteFloat32 at the current Offset
func (p *Packet) WriteFloat32(i float32) *Packet {
	binary.LittleEndian.PutUint32(p.Data[p.Offset:], math.Float32bits(i))
	p.Offset += (32 / 8)
	return p
}

// WriteInt64 at the current Offset
func (p *Packet) WriteInt64(i int64) *Packet {
	binary.LittleEndian.PutUint64(p.Data[p.Offset:], uint64(i))
	p.Offset += (64 / 8)
	return p
}

// WriteUInt64 at the current Offset
func (p *Packet) WriteUInt64(i uint64) *Packet {
	binary.LittleEndian.PutUint64(p.Data[p.Offset:], i)
	p.Offset += (64 / 8)
	return p
}

// WriteInt32 at the current Offset
func (p *Packet) WriteInt32(i int32) *Packet {
	binary.LittleEndian.PutUint32(p.Data[p.Offset:], uint32(i))
	p.Offset += (32 / 8)
	return p
}

// WriteUInt32 at the current Offset
func (p *Packet) WriteUInt32(i uint32) *Packet {
	binary.LittleEndian.PutUint32(p.Data[p.Offset:], i)
	p.Offset += (32 / 8)
	return p
}

// WriteInt16 at the current Offset
func (p *Packet) WriteInt16(i int16) *Packet {
	binary.LittleEndian.PutUint16(p.Data[p.Offset:], uint16(i))
	p.Offset += (16 / 8)
	return p
}

// WriteUInt16 at the current Offset
func (p *Packet) WriteUInt16(i uint16) *Packet {
	binary.LittleEndian.PutUint16(p.Data[p.Offset:], i)
	p.Offset += (16 / 8)
	return p
}

// WriteInt8 at the current Offset
func (p *Packet) WriteInt8(i int8) *Packet {
	p.Data[p.Offset] = uint8(i)
	p.Offset += (8 / 8)
	return p
}

// WriteUInt8 at the current Offset
func (p *Packet) WriteUInt8(i uint8) *Packet {
	p.Data[p.Offset] = i
	p.Offset += (8 / 8)
	return p
}

// Write3DVector at the current Offset
func (p *Packet) Write3DVector(vec *r3.Vector) *Packet {
	return p.WriteFloat32(float32(vec.X)).
		WriteFloat32(float32(vec.Y)).
		WriteFloat32(float32(vec.Z))
}

// WriteString (size+string) at the current Offset
func (p *Packet) WriteString(s string) *Packet {
	length := len(s)
	if length < 1 {
		return p
	}

	p.WriteInt32(int32(length))
	for i := 0; i < length; i++ {
		p = p.WriteUInt8(s[i])
	}

	return p
}

// ReadPacket create a new packet instance with the given input Data
func ReadPacket(d []byte) *Packet {
	p := new(Packet)
	p.Data = d

	return p
}

// ReadFloat32 at the current Offset
func (p *Packet) ReadFloat32() float32 {
	i := binary.LittleEndian.Uint32(p.Data[p.Offset:])
	p.Offset += (32 / 8)
	return math.Float32frombits(i)
}

// ReadInt64 at the current Offset
func (p *Packet) ReadInt64() int64 {
	i := binary.LittleEndian.Uint64(p.Data[p.Offset:])
	p.Offset += (64 / 8)
	return int64(i)
}

// ReadUInt64 at the current Offset
func (p *Packet) ReadUInt64() uint64 {
	i := binary.LittleEndian.Uint64(p.Data[p.Offset:])
	p.Offset += (64 / 8)
	return i
}

// ReadInt32 at the current Offset
func (p *Packet) ReadInt32() int32 {
	i := binary.LittleEndian.Uint32(p.Data[p.Offset:])
	p.Offset += (32 / 8)
	return int32(i)
}

// ReadUInt32 at the current Offset
func (p *Packet) ReadUInt32() uint32 {
	i := binary.LittleEndian.Uint32(p.Data[p.Offset:])
	p.Offset += (32 / 8)
	return i
}

// ReadInt16 at the current Offset
func (p *Packet) ReadInt16() int16 {
	i := binary.LittleEndian.Uint16(p.Data[p.Offset:])
	p.Offset += (16 / 8)
	return int16(i)
}

// ReadUInt16 at the current Offset
func (p *Packet) ReadUInt16() uint16 {
	i := binary.LittleEndian.Uint16(p.Data[p.Offset:])
	p.Offset += (16 / 8)
	return i
}

// ReadInt8 at the current Offset
func (p *Packet) ReadInt8() int8 {
	i := p.Data[p.Offset]
	p.Offset += (8 / 8)
	return int8(i)
}

// ReadUInt8 at the current Offset
func (p *Packet) ReadUInt8() uint8 {
	i := p.Data[p.Offset]
	p.Offset += (8 / 8)
	return i
}

// Read3DVector at the current Offset
func (p *Packet) Read3DVector() *r3.Vector {
	vec := new(r3.Vector)
	vec.X = float64(p.ReadFloat32())
	vec.Y = float64(p.ReadFloat32())
	vec.Z = float64(p.ReadFloat32())
	return vec
}

// ReadString (size+string) at the current Offset
func (p *Packet) ReadString() string {
	var buffer strings.Builder
	len := p.ReadUInt32()

	for i := uint32(0); i < len; i++ {
		buffer.WriteByte(p.ReadUInt8())
	}

	return buffer.String()
}
