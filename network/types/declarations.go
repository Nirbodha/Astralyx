package types

type (
	String string
	VariableInteger int32 //Why is this needed? Verifying the integrity of a packet, of course.
	Byte byte
	Float float32
	Integer int32
	ByteArray []Byte
)
