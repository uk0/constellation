package passwd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	filename := "/etc/passwd"

	testCases := map[string]struct {
		passwdContents  string
		createFile      bool
		expectedEntries Entries
		expectErr       bool
	}{
		"parse works": {
			passwdContents: "root:x:0:0:root:/root:/bin/bash\n",
			createFile:     true,
			expectedEntries: Entries{
				"root": {
					Pass:  "x",
					Uid:   "0",
					Gid:   "0",
					Gecos: "root",
					Home:  "/root",
					Shell: "/bin/bash",
				},
			},
			expectErr: false,
		},
		"passwd is corrupt": {
			passwdContents: "too:few:fields\n",
			createFile:     true,
			expectErr:      true,
		},
		"file does not exist": {
			createFile: false,
			expectErr:  true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			fs := afero.NewMemMapFs()
			if tc.createFile {
				assert.NoError(afero.WriteFile(fs, filename, []byte(tc.passwdContents), 0o644))
			}
			passwd := Passwd{}
			entries, err := passwd.Parse(fs)

			if tc.expectErr {
				assert.Error(err)
				return
			}
			require.NoError(err)
			assert.Equal(tc.expectedEntries, entries)
		})
	}
}
