package req


import "fmt"


/* HELPER FUNCTIONS ------------------------*/
// i imagine that these HELPER functions might
// be faster if they worked w/ indices
// instead of slices, since it probably
// requires less copying, but ill do it
// with slices to try them out better

// serialize a 16bit unsigned in network order
// returns slice to write to next
func SerU16(n uint16, tgt []byte) []byte {
	tgt[0] = byte(n >> 8)
	tgt[1] = byte(n)

	if len(tgt) < 3 {
		return nil
	}

	return tgt[2:]
}

// serialize a 32bit unsigned in network order
// returns slice to write to next
func SerU32(n uint32, tgt []byte) []byte {
	tgt[0] = byte(n >> 24)
	tgt[1] = byte(n >> 16)
	tgt[2] = byte(n >> 8)
	tgt[3] = byte(n)

	if len(tgt) < 5 {
		return nil
	}

	return tgt[4:]
}

// deserialize a 16bit unsigned in network order
// returns slice to read from next
func DeserU16(n *uint16, src []byte) []byte {
	*n = uint16(src[0]<<8) | uint16(src[1])

	if len(src) < 3 {
		return nil
	}

	return src[2:]
}

// deserialize a 32bit unsigned in network order
// returns slice to read from next
func DeserU32(n *uint32, src []byte) []byte {
	*n = uint32(src[0]<<24) | uint32(src[1]<<16)
	*n |= uint32(src[2]<<8) | uint32(src[3])

	if len(src) < 5 {
		return nil
	}

	return src[4:]
}

/*------------------------------------------*/


func ReqSerial(buf []byte, r *Req) error {

	if len(buf) < ReqBufSize {
		return fmt.Errorf("buf's too small to fit a Req dood")
	}

	writeh := buf[0:]
	writeh = SerU16(r.Id, writeh)
	writeh = SerU16(r.Rtype, writeh)
	writeh = SerU32(r.Src, writeh)
	writeh = SerU32(r.Info, writeh)
	writeh = SerU32(r.Plsz, writeh)

	return nil
}

func ReqDeserial(r *Req, buf []byte) error {

	if cap(buf) < ReqBufSize {
		return fmt.Errorf("buf's too small to have a Req dood")
	}

	readh := buf[0:]
	readh = DeserU16(&r.Id, readh)
	readh = DeserU16(&r.Rtype, readh)
	readh = DeserU32(&r.Src, readh)
	readh = DeserU32(&r.Info, readh)
	readh = DeserU32(&r.Plsz, readh)

	return nil
}

