package fcover

type Fcover map[string]struct{}

func (cov *Fcover) Merge(raw []string) {
	c := *cov
	if c == nil {
		c = make(Fcover)
		*cov = c
	}
	for _, pc := range raw {
		c[pc] = struct{}{}
	}
}

func (cov Fcover) Serialize() []string {
	res := make([]string, 0, len(cov))
	for pc := range cov {
		res = append(res, pc)
	}
	return res
}
