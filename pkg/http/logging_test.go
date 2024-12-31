package http

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_writeLog(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name    string
		db      *sql.DB
		log     log
		wantErr bool
	}{
		{
			name: "successful insert",
			db:   db,
			log: log{
				Epoch:   1234567890,
				Remote:  "192.168.1.1",
				Country: "US",
				Proto:   "HTTP/1.1",
				Host:    "example.com",
				Method:  "GET",
				Request: "/path",
				Status:  200,
				Bytes:   1024,
			},
			wantErr: false,
		},
		{
			name:    "nil db",
			db:      nil,
			log:     log{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &logWriter{
				db: tt.db,
			}

			if tt.db != nil {
				mock.ExpectExec("INSERT INTO http").
					WithArgs(tt.log.Epoch, tt.log.Remote, tt.log.Country, tt.log.Proto, tt.log.Host, tt.log.Method, tt.log.Request, tt.log.Status, tt.log.Bytes).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			w.writeLog(tt.log)

			if tt.db != nil {
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("unfulfilled expectations: %s", err)
				}
			}
		})
	}
}

func Test_loggingMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		remoteAddr string
		method     string
		host       string
		uri        string
		proto      string
		wantStatus int
	}{
		{
			name:       "successful request",
			remoteAddr: "192.168.1.1:1234",
			method:     "GET",
			host:       "example.com",
			uri:        "/test",
			proto:      "HTTP/1.1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid remote addr",
			remoteAddr: "invalid",
			method:     "POST",
			host:       "",
			uri:        "/path",
			proto:      "HTTP/1.1",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock db: %v", err)
			}
			defer db.Close()

			writer := &logWriter{db: db}

			mock.ExpectExec("INSERT INTO http").
				WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), tt.proto, sqlmock.AnyArg(), tt.method, sqlmock.AnyArg(), tt.wantStatus, sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 1))

			handler := loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.wantStatus)
			}), writer)

			req := httptest.NewRequest(tt.method, tt.uri, nil)
			req.RemoteAddr = tt.remoteAddr
			req.Host = tt.host
			req.Proto = tt.proto

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("loggingMiddleware() status = %v, want %v", rec.Code, tt.wantStatus)
			}

			time.Sleep(100 * time.Millisecond)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}
