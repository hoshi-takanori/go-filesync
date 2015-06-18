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
		SyncFile(mode, dir, f.Name, FInfo{}, f)
	}

	return result, nil
}

func SyncFile(mode int, dir Dir, name string, l, r FInfo) *FInfo {
	if l.ModTime.Unix() == r.ModTime.Unix() {
		if l.Size != r.Size {
			println("SIZE DIFFER: local", l.Size, "remote", r.Size)
		} else if l.Size == 0 {
			println("SIZE ZERO")
		}
		if mode == SyncModeSend {
			return &l
		}
	} else if l.ModTime.Unix() > r.ModTime.Unix() && l.Size > 0 {
		if mode == SyncModeSend || mode == SyncModeBoth {
			err := dir.Read(&l)
			if err != nil {
				println("read failed:", err.Error())
			}
			return &l
		}
	} else if l.ModTime.Unix() < r.ModTime.Unix() && r.Size > 0 {
		if mode == SyncModeBoth || mode == SyncModeReceive {
			err := dir.Write(r)
			if err != nil {
				println("write failed:", err.Error())
			}
		}
	} else if l.Name != "" {
		if mode == SyncModeBoth || mode == SyncModeReceive {
			err := dir.Remove(name)
			if err != nil {
				println("remove failed:", err.Error())
			}
		}
		if mode == SyncModeSend || mode == SyncModeBoth {
			return &l
		}
	}
	return nil
}
