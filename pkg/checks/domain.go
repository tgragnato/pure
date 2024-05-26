package checks

import "strings"

var whitelist = []string{
	"apple-finance.query.yahoo.com.",
	"blob.core.windows.net.",
	"d1612au5bhln8.cloudfront.net.",
	"dnfed9a4ru7oh.cloudfront.net.",
	"gcs-eu-00002.content-storage-download.googleapis.com.",
	"gcs-eu-00002.content-storage-upload.googleapis.com.",
	"semrush.com.",
}

func CheckDomain(domain string) bool {

	for _, allowed := range whitelist {
		if domain == allowed {
			return true
		}
	}

	for _, suffix := range blacklist {
		if strings.HasSuffix(domain, suffix) {
			return false
		}
	}

	return true
}
