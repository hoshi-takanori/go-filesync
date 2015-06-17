package main

const (
	SyncModeSend = iota
	SyncModeBoth
	SyncModeReceive
)

func SyncFiles(mode int, dir Dir, remoteFs []FInfo) ([]FInfo, error) {
	m := map[string]FInfo{}
	for _, f := range remoteFs {
		m[f.Name] = f
	}

	localFs, err := dir.List()
	if err != nil {
		return nil, err
	}

	result := []FInfo{}

	for _, f := range localFs {
		r := SyncFile(mode, dir, f.Name, f, m[f.Name])
		if r != nil {
			result = append(result, *r)
		}
		delete(m, f.Name)
	}

	for _, f := range m {
		r := SyncFile(mode, dir, f.Name, FInfo{}, f)
		if r != nil {
			result = append(result, *r)
		}
	}

	return result, nil
}

func SyncFile(mode int, dir Dir, name string, l, r FInfo) *FInfo {
	if l.ModTime.Unix() == r.ModTime.Unix() {
		println("skip", name)
		if l.Size != r.Size {
			println("SIZE DIFFER: local", l.Size, "remote", r.Size)
		} else if l.Size == 0 {
			println("SIZE ZERO")
		}
		if mode == SyncModeSend {
			return &l
		}
	} else if l.ModTime.Unix() > r.ModTime.Unix() && l.Size > 0 {
		println("put", name)
		if mode == SyncModeSend || mode == SyncModeBoth {
			dir.Read(&l)
			return &l
		}
	} else if l.ModTime.Unix() < r.ModTime.Unix() && r.Size > 0 {
		println("get", name)
		dir.Write(r)
		if l.Name != "" && mode == SyncModeSend {
			return &l
		}
	} else {
		if l.Name != "" {
			println("rm local", name)
			dir.Remove(name)
		}
		if r.Name != "" {
			println("rm remote", name)
			if l.Name != "" && (mode == SyncModeSend || mode == SyncModeBoth) {
				return &l
			}
		}
	}
	return nil
}
