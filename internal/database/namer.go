package database

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// AS_ADDR_OBJ_20250626_bc6f64d9-fb28-40d6-8a99-57e44b920d07.XML => addr_obj
// AS_CHANGE_HISTORY_20250626_d1a57485-156c-4463-8a23-2328fb0f6f9d => change_history
func ResolveTableName(filename string) (string, error) {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))

	re := regexp.MustCompile(`(?i)^AS_([A-Z0-9_]+?)_\d{8}`)
	if m := re.FindStringSubmatch(base); len(m) == 2 {
		return strings.ToLower(m[1]), nil
	}
	return "", fmt.Errorf("cannot resolve table name from filename: %s", filename)
}

// ITEM_ID => item_id
// CHANGEID => changeid
func ResolveColumnName(field string) string {
	return strings.ToLower(field)
}
