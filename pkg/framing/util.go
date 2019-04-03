package framing

func hasDictinctFrames(fs *[]Frame) bool {
	if len(*fs) < 2 {
		return false
	}

	fst := (*fs)[0]
	df := MetaDataString(&(fst.MetaData))

	// return true the first time you see
	// a different frame
	for _, f := range (*fs)[1:] {
		mds := MetaDataString(&(f.MetaData))
		if df != mds {
			return true
		}
	}

	return false
}

func inSlice(sl *[]string, s string) bool {
	for _, v := range *sl {
		if s == v {
			return true
		}
	}
	return false
}

func appendIfNew(sl *[]string, s string) *[]string {
	var nsl []string
	if !inSlice(sl, s) {
		nsl = append(*sl, s)
		return &nsl
	}
	return sl
}

func concatUnique(sl *[]string, sl2 *[]string) *[]string {
	for _, v := range *sl2 {
		sl = appendIfNew(sl, v)
	}
	return sl
}
