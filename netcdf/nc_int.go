// Copyright 2014 The Go-NetCDF Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// These files are autogenerated from nc_double.go using generate.go
// DO NOT EDIT (except nc_double.go).

package netcdf

import (
	"fmt"
	"unsafe"
)

// #include <stdlib.h>
// #include <netcdf.h>
import "C"

// WriteInt32s writes data as the entire data for variable v.
func (v Var) WriteInt32s(data []int32) error {
	if err := okData(v, INT, len(data)); err != nil {
		return err
	}
	return newError(C.nc_put_var_int(C.int(v.ds), C.int(v.id), (*C.int)(unsafe.Pointer(&data[0]))))
}

// ReadInt32s reads the entire variable v into data, which must have enough
// space for all the values (i.e. len(data) must be at least v.Len()).
func (v Var) ReadInt32s(data []int32) error {
	if err := okData(v, INT, len(data)); err != nil {
		return err
	}
	return newError(C.nc_get_var_int(C.int(v.ds), C.int(v.id), (*C.int)(unsafe.Pointer(&data[0]))))
}

func (v Var) ReadArrayInt32s(offsets []int, lens []int, data []int32) error {
	dims, err := v.Dims()
	if err != nil {
		return err
	}
	if len(offsets) != len(dims) || len(lens) != len(dims) {
		return fmt.Errorf("Invalid array of offsets/lens")
	}

	starts := make([]C.size_t, len(dims))
	counts := make([]C.size_t, len(dims))
	total_len := 1
	for i, d := range dims {
		dim_len, _ := d.Len()
		if lens[i] == -1 {
			lens[i] = int(dim_len) - offsets[i]
		}
		if offsets[i]+lens[i] > int(dim_len) {
			return fmt.Errorf("Invalid offset/len %d >= %d", i, dim_len)
		}
		starts[i] = C.size_t(offsets[i])
		counts[i] = C.size_t(lens[i])
		total_len = total_len * lens[i]
	}

	u, err := v.Type()
	if err != nil {
		return err
	}
	if u != INT {
		return fmt.Errorf("wrong data type %s; expected %s", u, INT)
	}
	if len(data) != total_len {
		return fmt.Errorf("Invalid length of recieving data %d (need %d)", len(data), total_len)
	}

	return newError(C.nc_get_vara_int(
		C.int(v.ds), C.int(v.id),
		(*C.size_t)(unsafe.Pointer(&starts[0])),
		(*C.size_t)(unsafe.Pointer(&counts[0])),
		(*C.int)(unsafe.Pointer(&data[0]))))
}

// WriteInt32s sets the value of attribute a to val.
func (a Attr) WriteInt32s(val []int32) error {
	// We don't need okData here because netcdf library doesn't know
	// the length or type of the attribute yet.
	cname := C.CString(a.name)
	defer C.free(unsafe.Pointer(cname))
	return newError(C.nc_put_att_int(C.int(a.v.ds), C.int(a.v.id), cname,
		C.nc_type(INT), C.size_t(len(val)), (*C.int)(unsafe.Pointer(&val[0]))))
}

// ReadInt32s reads the entire attribute value into val.
func (a Attr) ReadInt32s(val []int32) (err error) {
	if err := okData(a, INT, len(val)); err != nil {
		return err
	}
	cname := C.CString(a.name)
	defer C.free(unsafe.Pointer(cname))
	err = newError(C.nc_get_att_int(C.int(a.v.ds), C.int(a.v.id), cname,
		(*C.int)(unsafe.Pointer(&val[0]))))
	return
}

// Int32sReader is a interface that allows reading a sequence of values of fixed length.
type Int32sReader interface {
	Len() (n uint64, err error)
	ReadInt32s(val []int32) (err error)
}

// GetInt32s reads the entire data in r and returns it.
func GetInt32s(r Int32sReader) (data []int32, err error) {
	n, err := r.Len()
	if err != nil {
		return
	}
	data = make([]int32, n)
	err = r.ReadInt32s(data)
	return
}

// testReadInt32s writes somes data to v. N is v.Len().
// This function is only used for testing.
func testWriteInt32s(v Var, n uint64) error {
	data := make([]int32, n)
	for i := 0; i < int(n); i++ {
		data[i] = int32(i + 10)
	}
	return v.WriteInt32s(data)
}

// testReadInt32s reads data from v and checks that it's the same as what
// was written by testWriteDouble. N is v.Len().
// This function is only used for testing.
func testReadInt32s(v Var, n uint64) error {
	data := make([]int32, n)
	if err := v.ReadInt32s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if val := int32(i + 10); data[i] != val {
			return fmt.Errorf("data at position %d is %v; expected %v\n", i, data[i], val)
		}
	}

	data = data[:12]
	if err := v.ReadArrayInt32s([]int{1, 1}, []int{-1, -1}, data); err != nil {
		return err
	}
	if data[0] != int32(14) {
		return fmt.Errorf("data as sub-array[0] != 14")
	}
	if data[11] != int32(30) {
		return fmt.Errorf("data as sub-array[11] != 30")
	}
	// fmt.Printf("array float64 %d %s\n", len(data), data)
	return nil
}
