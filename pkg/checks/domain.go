package checks

import "strings"

var whitelist = []string{
	"apple-finance.query.yahoo.com.",
	"d3w4wp6rol9nvz.cloudfront.net.",
	"dnfed9a4ru7oh.cloudfront.net.",
	"gcs-eu-00002.content-storage-download.googleapis.com.",
	"gcs-eu-00002.content-storage-upload.googleapis.com.",
	"mattermost.pepita.io.",
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