package framing

type collectFunc func(fr Frame, cmap map[string]bool) bool

func collect(f collectFunc, fs *[]Frame) []*Frame {
	tmp := make(map[string]bool)
	cfs := []*Frame{}
	// do not use pointers from values returned from range
	for i, fr := range *fs {
		mds := MetaDataString(&fr.MetaData)
		if f(fr, tmp) == true {
			tmp[mds] = true
			cfs = append(cfs, &((*fs)[i]))
		}
	}

	return cfs
}

func collectDistinct(f Frame, cmap map[string]bool) bool {
	mds := MetaDataString(&f.MetaData)
	if _, v := cmap[mds]; !v {
		return true
	}
	return false
}
