package main

import ()

func SyncFiles(dir Dir, remoteFs []FInfo) ([]FInfo, error) {
	m := map[string]FInfo{}
	for _, f := range remoteFs {
		m[f.Name] = f
	}

	localFs, err := dir.List()
	if err != nil {
		return []FInfo{}, err
	}

	for _, f := range localFs {
		SyncFile(dir, f.Name, f, m[f.Name])
		delete(m, f.Name)
	}

	for _, f := range m {
		SyncFile(dir, f.Name, FInfo{}, f)
	}

	return []FInfo{}, nil
}

func SyncFile(dir Dir, name string, l, r FInfo) {
	if l.ModTime.Unix() == r.ModTime.Unix() {
		println("skip", name)
		if l.Size != r.Size {
			println("SIZE DIFFER: local", l.Size, "remote", r.Size)
		} else if l.Size == 0 {
			println("SIZE ZERO")
		}
	} else if l.ModTime.Unix() > r.ModTime.Unix() && l.Size > 0 {
		println("put", name)
	} else if l.ModTime.Unix() < r.ModTime.Unix() && r.Size > 0 {
		println("get", name)
		dir.Write(r)
	} else {
		if l.Name != "" {
			println("rm local", name)
			dir.Remove(name)
		}
		if r.Name != "" {
			println("rm remote", name)
		}
	}
}
