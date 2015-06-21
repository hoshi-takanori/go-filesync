package main

const (
	SyncModeBegin = iota
	SyncModeSend
	SyncModeBoth
	SyncModeReceive
)

func SyncFiles(mode int, dir Dir, remoteFs []FInfo) ([]FInfo, error) {
	dir.Log("list " + dir.Path(""))
	localFs, err := dir.List()
	if mode == SyncModeBegin || err != nil {
		return localFs, err
	}

	result := []FInfo{}

	m := map[string]FInfo{}
	for _, f := range remoteFs {
		m[f.Name] = f
	}

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
			dir.Log("size differ " + dir.Path(name))
		} else if l.Size == 0 {
			dir.Log("size zero " + dir.Path(name))
		}
		if mode == SyncModeSend {
			return &l
		}
	} else if l.ModTime.Unix() > r.ModTime.Unix() && l.Size > 0 {
		if mode == SyncModeSend || mode == SyncModeBoth {
			dir.Log("read " + dir.Path(name))
			err := dir.Read(&l)
			if err != nil {
				dir.Log("read failed: " + err.Error())
			}
			return &l
		}
	} else if l.ModTime.Unix() < r.ModTime.Unix() && r.Size > 0 {
		if mode == SyncModeBoth || mode == SyncModeReceive {
			dir.Log("save " + dir.Path(name))
			err := dir.Write(r)
			if err != nil {
				dir.Log("write failed: " + err.Error())
			}
		}
	} else if l.Name != "" {
		if mode == SyncModeBoth || mode == SyncModeReceive {
			dir.Log("rm " + dir.Path(name))
			err := dir.Remove(name)
			if err != nil {
				dir.Log("remove failed: " + err.Error())
			}
		}
		if mode == SyncModeSend || mode == SyncModeBoth {
			return &l
		}
	}
	return nil
}
