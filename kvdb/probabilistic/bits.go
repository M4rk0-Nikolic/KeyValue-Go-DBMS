package bits

func ChangeNthBitTo(array *[]uint64, n uint, flag uint) { // flag  0 to set to 0, 1 to set to 1, 2 to flip
	len := array.Len()
	block64offset = ((len + 63) >> 6) - 1

	mask := uint64(1 << (n % 64))

	switch flag{
	case 0:
		(*array)[block64offset] = (*array)[block64offset] &^ mask
	case 1:
		(*array)[block64offset] = (*array)[block64offset] | mask
	case 2:
		(*array)[block64offset] = (*array)[block64offset] ^ mask
	default:
		return
	}
}

func IsNthBitOne(array *[]uint64, n uint){
	len := array.Len()
	block64offset = ((len + 63) >> 6) - 1

	mask := uint64(1 << (n % 64))

	if (((*array)[block64offset] & mask) != 0 ){
		return True
	}
}

func NumberOfTrailingZeros(){
	// for hyper log log
}
