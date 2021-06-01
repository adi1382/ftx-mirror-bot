package optional

type Float64 struct {
	set   bool
	value float64
}

func (o *Float64) IsSet() bool {
	return o.set
}

func (o *Float64) Value() float64 {
	return o.value
}

func (o *Float64) Set(v float64) {
	o.set = true
	o.value = v
}

type Int64 struct {
	set   bool
	value int64
}

func (o *Int64) IsSet() bool {
	return o.set
}

func (o *Int64) Value() int64 {
	return o.value
}

func (o *Int64) Set(v int64) {
	o.set = true
	o.value = v
}

type String struct {
	set   bool
	value string
}

func (o *String) IsSet() bool {
	return o.set
}

func (o *String) Value() string {
	return o.value
}

func (o *String) Set(v string) {
	o.set = true
	o.value = v
}

type Bool struct {
	set   bool
	value bool
}

func (o *Bool) IsSet() bool {
	return o.set
}

func (o *Bool) Value() bool {
	return o.value
}

func (o *Bool) Set(v bool) {
	o.set = true
	o.value = v
}
