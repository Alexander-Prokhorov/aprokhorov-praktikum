package files_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"aprokhorov-praktikum/internal/server/files"
	"aprokhorov-praktikum/internal/storage"
)

func TestSaveData(t *testing.T) {
	t.Parallel()

	type args struct {
		fileName string
		s        storage.Storage
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test for writing data to file",
			args: args{
				fileName: "testdata/test.txt",
				s:        storage.NewStorageMem(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := files.SaveData(tt.args.fileName, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("SaveData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadData(t *testing.T) {
	t.Parallel()

	type args struct {
		fileName string
		s        storage.Storage
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test for loading data from file",
			args: args{
				fileName: "testdata/test.txt",
				s:        storage.NewStorageMem(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			er := files.SaveData(tt.args.fileName, tt.args.s)
			assert.NoError(t, er, "SaveData to File Error")
			if err := files.LoadData(tt.args.fileName, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("LoadData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
