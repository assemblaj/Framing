package framing

type groupFunc func(fr Frame, gmap map[string][]*Frame) (bool, string)

func group(f groupFunc, fs *[]Frame) map[string][]*Frame {
	gfs := make(map[string][]*Frame)
	for i, fr := range *fs {
		b, k := f(fr, gfs)
		if b == true {
			gfs[k] = append(gfs[k], &((*fs)[i]))
		}
	}
	return gfs
}

func groupMetadata(fr Frame, gmap map[string][]*Frame) (bool, string) {
	mds := MetaDataString(&fr.MetaData)
	return true, mds
}

func hasGroups(fs *[]Frame) bool {
	mdsl := []string{}
	for _, f := range *fs {
		mds := MetaDataString(&(f.MetaData))
		if inSlice(&mdsl, mds) {
			return true
		}
		mdsl = append(mdsl, mds)
	}
	return false
}
