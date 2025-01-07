package filetrove

import "github.com/pkg/xattr"

// GetXattr checks if an inpde has xattr. It returns a list of names and values.
func GetXattr(filePath string) (map[string]string, error) {
	var list []string
	filexattrmap := make(map[string]string)

	list, err := xattr.List(filePath)
	if err != nil {
		return filexattrmap, err
	}

	for _, v := range list {
		xattrbyte, err := xattr.Get(filePath, v)
		if err != nil {
			return filexattrmap, err
		}
		filexattrmap[v] = string(xattrbyte)
	}
	return filexattrmap, err
}
