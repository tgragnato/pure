package dnshandlers

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCrossPrefetch(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	d := &DnsHandlers{
		db: db,
	}

	tests := []struct {
		name    string
		rows    *sqlmock.Rows
		wantErr bool
	}{
		{
			name: "successful query",
			rows: sqlmock.NewRows([]string{"key"}).
				AddRow("example.com"),
			wantErr: false,
		},
		{
			name:    "query error",
			rows:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				mock.ExpectQuery("SELECT a.key").WillReturnError(sql.ErrConnDone)
			} else {
				mock.ExpectQuery("SELECT a.key").WillReturnRows(tt.rows)
			}

			d.crossPrefetch()

			err := mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("expectations not met: %v", err)
			}
		})
	}
}

func TestSelfPrefetch(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	d := &DnsHandlers{
		db: db,
	}

	tests := []struct {
		name    string
		ipv6    bool
		muIPv6  bool
		muIPv4  bool
		rows    *sqlmock.Rows
		wantErr bool
	}{
		{
			name:    "successful ipv4 query",
			ipv6:    false,
			muIPv4:  false,
			rows:    sqlmock.NewRows([]string{"key"}).AddRow("example.com"),
			wantErr: false,
		},
		{
			name:    "successful ipv6 query",
			ipv6:    true,
			muIPv6:  false,
			rows:    sqlmock.NewRows([]string{"key"}).AddRow("example.com"),
			wantErr: false,
		},
		{
			name:    "mutex locked ipv4",
			ipv6:    false,
			muIPv4:  true,
			wantErr: false,
		},
		{
			name:    "mutex locked ipv6",
			ipv6:    true,
			muIPv6:  true,
			wantErr: false,
		},
		{
			name:    "query error",
			ipv6:    false,
			muIPv4:  false,
			rows:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d.muIPv6 = tt.muIPv6
			d.muIPv4 = tt.muIPv4

			if !tt.muIPv4 && !tt.ipv6 || !tt.muIPv6 && tt.ipv6 {
				if tt.wantErr {
					mock.ExpectQuery("SELECT key").WillReturnError(sql.ErrConnDone)
				} else {
					mock.ExpectQuery("SELECT key").WillReturnRows(tt.rows)
				}
			}

			d.selfPrefetch(tt.ipv6)

			err := mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("expectations not met: %v", err)
			}
		})
	}
}
