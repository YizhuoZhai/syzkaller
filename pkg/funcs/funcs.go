package funcs
// Package signal provides types for working with feedback functions.
type (
	elemType string
	prioType int8
)

type Funcs map[elemType]prioType

type Serial struct {
	Elems []elemType
	Prios []prioType
}

func (f Funcs) Len() int {
	return len(f)
}

func (f Funcs) Empty() bool {
	return len(f) == 0
}
func (f Funcs) Copy() Funcs {
	c := make(Funcs, len(f))
	for e, p := range f {
		c[e] = p
	}
	return c
}

func (f *Funcs) Split(n int) Funcs {
	if f.Empty() {
		return nil
	}
	c := make(Funcs, n)
	for e, p := range *f {
		delete(*f, e)
		c[e] = p
		n--
		if n == 0 {
			break
		}
	}
	if len(*f) == 0 {
		*f = nil
	}
	return c
}

func FromRaw(raw []string, prio uint8) Funcs {
	if len(raw) == 0 {
		return nil
	}
	s := make(Funcs, len(raw))
	for _, e := range raw {
		//set the priority for every function into prio
		s[elemType(e)] = prioType(prio)
	}
	return s
}

func (f Funcs) Serialize() Serial {
	if f.Empty() {
		return Serial{}
	}
	res := Serial{
		Elems: make([]elemType, len(f)),
		Prios: make([]prioType, len(f)),
	}
	i := 0
	for e, p := range f {
		res.Elems[i] = e
		res.Prios[i] = p
		i++
	}
	return res
}
func (ser Serial) Deserialize() Funcs {
	if len(ser.Elems) != len(ser.Prios) {
		panic("corrupted Serial")
	}
	if len(ser.Elems) == 0 {
		return nil
	}
	s := make(Funcs, len(ser.Elems))
	for i, e := range ser.Elems {
		s[e] = ser.Prios[i]
	}
	return s
}

func (f Funcs) Diff(f1 Funcs) Funcs {
	if f1.Empty() {
		return nil
	}
	var res Funcs
	for e, p1 := range f1 {
		if p, ok := f[e]; ok && p >= p1 {
			continue
		}
		if res == nil {
			res = make(Funcs)
		}
		res[e] = p1
	}
	return res
}
func (f Funcs) DiffRaw(raw []string, prio uint8) Funcs {
	var res Funcs
	for _, e := range raw {
		if p, ok := f[elemType(e)]; ok && p >= prioType(prio) {
			continue
		}
		if res == nil {
			res = make(Funcs)
		}
		res[elemType(e)] = prioType(prio)
	}
	return res
}

func (f Funcs) Intersection(f1 Funcs) Funcs {
	if f1.Empty() {
		return nil
	}
	res := make(Funcs, len(f))
	for e, p := range f {
		if p1, ok := f1[e]; ok && p1 >= p {
			res[e] = p
		}
	}
	return res
}
func (f *Funcs) Merge(f1 Funcs) {
	if f1.Empty() {
		return
	}
	s0 := *f
	if s0 == nil {
		s0 = make(Funcs, len(f1))
		*f = s0
	}
	for e, p1 := range f1 {
		if p, ok := s0[e]; !ok || p < p1 {
			s0[e] = p1
		}
	}
}
type Context struct {
	Funcs  Funcs
	Context interface{}
}

func Minimize(corpus []Context) []interface{} {
	type ContextPrio struct {
		prio prioType
		idx  int
	}
	covered := make(map[elemType]ContextPrio)
	for i, inp := range corpus {
		for e, p := range inp.Funcs {
			if prev, ok := covered[e]; !ok || p > prev.prio {
				covered[e] = ContextPrio{
					prio: p,
					idx:  i,
				}
			}
		}
	}
	indices := make(map[int]struct{}, len(corpus))
	for _, cp := range covered {
		indices[cp.idx] = struct{}{}
	}
	result := make([]interface{}, 0, len(indices))
	for idx := range indices {
		result = append(result, corpus[idx].Context)
	}
	return result
}