package checks

import "strings"

var whitelist = []string{
	"apple-finance.query.yahoo.com.",
	"d1612au5bhln8.cloudfront.net.",
	"d179pbpxvabqy1.cloudfront.net.",
	"d2gwrxj8kvy81r.cloudfront.net.",
	"d2ziyq7zksuzai.cloudfront.net.",
	"d36jcksde1wxzq.cloudfront.net.",
	"d3c9ouasuy8pg6.cloudfront.net.",
	"dcoo28a3i62o.cloudfront.net.",
	"dnfed9a4ru7oh.cloudfront.net.",
	"gcs-eu-00002.content-storage-download.googleapis.com.",
	"gcs-eu-00002.content-storage-upload.googleapis.com.",
}

var whitelistSuffix = []string{
	".azurewebsites.windows.net.",
	".blob.core.windows.net.",
	".blz25prdstrz09a.store.core.windows.net.",
	".blz25prdstrz09a.trafficmanager.net.",
	".bn9prdstrz04a.store.core.windows.net.",
	".bn9prdstrz04a.trafficmanager.net.",
	".ows.farm.",
	".semrush.com.",
}

func CheckDomain(domain string) bool {
	if !strings.HasSuffix(domain, ".") ||
		strings.Count(domain, ".") < 2 {
		return false
	}

	if strings.HasSuffix(domain, ".tgragnato.it.") &&
		domain != "www.tgragnato.it." &&
		domain != "stun.tgragnato.it." &&
		domain != "dht.tgragnato.it." &&
		domain != "api.tgragnato.it." {
		return false
	}

	for _, allowed := range whitelist {
		if domain == allowed {
			return true
		}
	}

	for _, suffix := range whitelistSuffix {
		if strings.HasSuffix(domain, suffix) {
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
