//go:build windows

package filetrove

import (
	"golang.org/x/sys/windows"
)

// getFileOwner gets owner and group on Windows
func GetFileOwner(path string) (*FileOwnerInfo, error) {
	fileName, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}

	handle, err := windows.CreateFile(fileName, windows.FILE_READ_ATTRIBUTES, windows.FILE_SHARE_READ, nil, windows.OPEN_EXISTING, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(handle)

	var secInfo windows.SECURITY_INFORMATION = windows.OWNER_SECURITY_INFORMATION | windows.GROUP_SECURITY_INFORMATION
	var secDesc *windows.SECURITY_DESCRIPTOR

	err = windows.GetSecurityInfo(handle, windows.SE_FILE_OBJECT, secInfo, nil, nil, nil, nil, &secDesc)
	if err != nil {
		return nil, err
	}

	// get owner
	var ownerSID *windows.SID
	err = windows.GetSecurityDescriptorOwner(secDesc, &ownerSID, nil)
	if err != nil {
		return nil, err
	}

	ownerName, err := lookupAccountName(ownerSID)
	if err != nil {
		return nil, err
	}

	// get group
	var groupSID *windows.SID
	err = windows.GetSecurityDescriptorGroup(secDesc, &groupSID, nil)
	if err != nil {
		return nil, err
	}

	groupName, err := lookupAccountName(groupSID)
	if err != nil {
		return nil, err
	}

	return &FileOwnerInfo{
		Owner: ownerName,
		Group: groupName,
	}, nil
}
