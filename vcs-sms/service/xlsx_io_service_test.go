package service

import (
	"io/fs"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportXLSX(t *testing.T) {
	type input struct {
		filePath string
	}
	type expectation struct {
		rows [][]string
		err  error
	}
	tests := map[string]struct {
		input
		expectation
	}{
		"Test_1": {
			input: input{
				filePath: "../tmp/24-05-04_11-37-17_template.xlsx",
			},
			expectation: expectation{
				rows: [][]string{},
				err:  nil,
			},
		},
		"Test_2": {
			input: input{
				filePath: "../tmp/24-05-04_11-37-17_template.xlsxxx",
			},
			expectation: expectation{
				rows: [][]string{},
				err: &fs.PathError{
					Op:   "open",
					Path: "../tmp/24-05-04_11-37-17_template.xlsxxx",
					Err:  syscall.Errno(2),
				},
			},
		},
		"Test_3_Sheet_No_Exist": {
			input: input{
				filePath: "../tmp/24-05-04_11-33-43_template.xlsx",
			},
			expectation: expectation{
				rows: [][]string{},
				err:  nil,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			service := NewXLSXService()
			_, err := service.ImportXLSX(test.input.filePath)
			assert.Equal(t, test.expectation.err, err)
		})
	}
}
