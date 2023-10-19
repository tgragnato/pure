package main

import (
	"strings"
)

func checkDomain(domain string) bool {
	for _, suffix := range blacklist {
		if strings.HasSuffix(domain, suffix) {
			return false
		}
	}

	return true
}

var blacklist = []string{
	// Reserved
	"arpa.",
	// China TLD
	"cn.",
	"xn--fiqs8s.",
	"xn--fiqz9s.",
	"hk.",
	"xn--j6w193g.",
	"mo.",
	"xn--mix082f.",
	"xn--mix891f.",
	// China domains
	"163.com.",
	"58.com.",
	"aliapp.org.",
	"alibabagroup.com.",
	"alicdn.com.",
	"aliexpress.com.",
	"alipay.com.",
	"alitrip.com.",
	"aliyuncs.com.",
	"amap.com.",
	"baidu.com.",
	"bilibili.com.",
	"byted.org.",
	"bytefcdn-oversea.com.",
	"byteoversea.com.",
	"csdn.net.",
	"ctrip.com.",
	"dbank.com.",
	"dbankcdn.com.",
	"dbankcloud.com.",
	"dbankcloud.eu.",
	"dbankedge.net.",
	"dianping.com.",
	"harmonyos.com.",
	"hicloud.com.",
	"huanqiu.com.",
	"huawei.asia.",
	"huawei.com.",
	"huaweistatic.com.",
	"hwccpc.com.",
	"ibytedtos.com.",
	"ibyteimg.com.",
	"ifeng.com.",
	"iqiyi.com.",
	"jd.com.",
	"meituan.com.",
	"mgtv.com.",
	"mi-img.com.",
	"mi.com.",
	"miui.com.",
	"mmstat.com.",
	"muscdn.com.",
	"musical.ly.",
	"nuomi.com.",
	"qingting.fm.",
	"qq.com.",
	"qunar.com.",
	"qyer.com.",
	"servicewechat.com.",
	"sogou.com.",
	"sohu.com.",
	"taobao.com.",
	"tiktok.com.",
	"tiktok.net.",
	"tiktokcdn.com.",
	"tiktokv.com.",
	"tmall.com.",
	"vmall.com.",
	"wechat.com.",
	"wechat.org.",
	"wechatapp.com.",
	"weibo.com.",
	"xiachufang.com.",
	"xiami.com.",
	"xiaomi.com.",
	"xiaomi.net.",
	"xiaomiyoupin.com.",
	"ximalaya.com.",
	"xinhuanet.com.",
	"yinyuetai.com.",
	"youku.com.",
	"zhihu.com.",
	// Russia TLD
	"ru.",
	"su.",
	"xn--p1ai.",
	"by.",
	"xn--90ais.",
	// Russia domains
	"evraz.com.",
	"kas-labs.com.",
	"kaspersky-labs.com.",
	"kaspersky.com.",
	"labkas.com.",
	"metalloinvest.com.",
	"n1mk.com.",
	"nangs.org.",
	"nornickel.com.",
	"polymetalinternational.com.",
	"rt.com.",
	"sverstal.com.",
	"uralkali.com.",
	"userapi.com.",
	"vk-cdn.net.",
	"vk-portal.net.",
	"vk.com.",
	"vk.company.",
	"vkuser.net.",
	"vkuseraudio.net.",
	"vkuservideo.net.",
	"yandex.",
	"yandex.com.",
	"yandex.eu.",
	"yandex.net.",
	// Iran TLD
	"ir.",
	"xn--mgba3a4f16a.",
	// Iran domains
	"agah.com.",
	"aparat.com.",
	"bale.ai.",
	"digikala.com.",
	"digimovie.vip.",
	"digimoviez.com.",
	"eitaa.com.",
	"emofid.com.",
	"filimo.com.",
	"mediaad.org.",
	"najva.com.",
	"rahavard365.com.",
	"sanjagh.com.",
	"sunista.info.",
	"telewebion.com.",
	"torob.com.",
	"tsetmc.com.",
	"varzesh3.com.",
	"yektanet.com.",
	// North Korea
	"kp.",
	// Blocks
	"1rx.io.",
	"360yield.com.",
	"6bgaput9ullc.com.",
	"ad-delivery.net.",
	"adcolony.com.",
	"addresseepaper.com.",
	"addthis.com.",
	"addtoany.com.",
	"adform.com.",
	"adform.net.",
	"adformdsp.net.",
	"adj.st.",
	"adjust.com.",
	"adjust.net.in.",
	"adjust.world.",
	"adkaora.space.",
	"adnxs-simple.com.",
	"adnxs.com.",
	"ads-twitter.com.",
	"adsafeprotected.com.",
	"adsco.re.",
	"adsensecustomsearchads.com.",
	"adsmoloco.com.",
	"adsrvr.org.",
	"advertising.com.",
	"advgo.net.",
	"amazon-adsystem.com.",
	"amplitude.com.",
	"ampproject.net.",
	"ampproject.org.",
	"analytics.edgekey.net.",
	"app-measurement.com.",
	"app-measurement.com/sdk-exp.",
	"appboy.eu.",
	"applovin.com.",
	"applvn.com.",
	"appsflyer.com.",
	"appsflyersdk.com.",
	"asbgdfxrau.com.",
	"banquetunarmedgrater.com.",
	"bodis.com.",
	"bounceexchange.com.",
	"brandmetrics.com.",
	"braze-images.com.",
	"braze.com.",
	"braze.eu.",
	"brightcove.com.",
	"brightcove.net.",
	"browser-intake-datadoghq.com.",
	"bugsnag.com.",
	"casalemedia.com.",
	"cbjqbkacuxjw.com.",
	"cdn4ads.com.",
	"cedexis-radar.net.",
	"cedexis.com.",
	"chartbeat.com.",
	"chartbeat.net.",
	"cloudflareinsights.com.",
	"commento.io.",
	"commentsmodule.com.",
	"comscore.com.",
	"comscoreresearch.com.",
	"concert.io.",
	"conde.digital.",
	"confiant-integrations.net.",
	"connexity.com.",
	"connexity.net.",
	"convertexperiments.com.",
	"cookiebot.com.",
	"cookielaw.org.",
	"crashlytics.com.",
	"crashlyticsreports-pa.googleapis.com.",
	"creative-bars1.com.",
	"criteo.com.",
	"criteo.net.",
	"crowdsignal.com.",
	"czyasezpvs.com.",
	"definedlaunching.com.",
	"demdex.net.",
	"discoveryfeed.org.",
	"displayvertising.com.",
	"disqus.com.",
	"disquscdn.com.",
	"dns.google.",
	"dns.google.com.",
	"dotmetrics.net.",
	"doubleclick.net.",
	"dynamicyield.com.",
	"eslzcjnlkepoow.com.",
	"evolutionadv.it.",
	"example.org.",
	"ezoic.com.",
	"ezoic.net.",
	"facebook.com.",
	"facebook.net.",
	"fastly-insights.com.",
	"fedfceddu.com.",
	"firebaselogging-pa.googleapis.com.",
	"fontawesome.com.",
	"fonts.googleapis.com.",
	"fonts.gstatic.com.",
	"fqdwrgbbkmlbh.com.",
	"friendshipmale.com.",
	"ftjcfx.com.",
	"fundingchoicesmessages.google.com.",
	"gatekeeperconsent.com.",
	"glyph.medium.com.",
	"gmxvmvptfm.com.",
	"google-analytics.com.",
	"googleadservices.com.",
	"googlesyndication.com.",
	"googletagmanager.com.",
	"googletagservices.com.",
	"gumgum.com.",
	"histats.com.",
	"hsadspixel.net.",
	"iinnbewhna.com.",
	"imasdk.googleapis.com.",
	"img-taboola.com.",
	"imrworldwide.com.",
	"inner-active.mobi.",
	"instagram.com.",
	"insurads.com.",
	"intake-analytics.wikimedia.org.",
	"intellipopup.com.",
	"intercom.io.",
	"intercomcdn.com.",
	"ipv4only.arpa.",
	"iubenda.com.",
	"jpyxqeysh.com.",
	"kustomerapp.com.",
	"lasubqueries.com.",
	"launchdarkly.com.",
	"logrocket-assets.io.",
	"ltmsphrcl.net.",
	"lunalabs.io.",
	"mambasms.com.",
	"marketo.net.",
	"matterlytics.com.",
	"memoinsights.com.",
	"moatads.com.",
	"mpi-internal.com.",
	"mpianalytics.com.",
	"nativery.com.",
	"newrelic.com.",
	"nocookie.net.",
	"nr-data.net.",
	"ojuhjcmhemvs.com.",
	"omnitagjs.com.",
	"onetag-sys.com.",
	"onetrust.com.",
	"onetrust.io.",
	"online-metrix.net.",
	"optimizely.com.",
	"outbrain.com.",
	"outbrainimg.com.",
	"padfungusunless.com.",
	"parse.ly.",
	"parsely.com.",
	"perfectmarket.com.",
	"permutive.app.",
	"permutive.com.",
	"piano.io.",
	"plausible.io.",
	"polldaddy.com.",
	"privacy-center.org.",
	"privacy-mgmt.com.",
	"proxy-safebrowsing.googleapis.com.",
	"pub.network.",
	"pubmatic.com.",
	"pubmnet.com.",
	"pubnative.net.",
	"pubtech.ai.",
	"pushio.com.",
	"quantcount.com.",
	"quantserve.com.",
	"raddoppia-bitcoin.click.",
	"rqmob.com.",
	"rubiconproject.com.",
	"run-syndicate.com.",
	"runative-syndicate.com.",
	"sail-horizon.com.",
	"sascdn.com.",
	"sb89347.com.",
	"schip.io.",
	"scorecardresearch.com.",
	"seadform.net.",
	"seedtag.com.",
	"sentry.io.",
	"servedby-buysellads.com.",
	"servenobid.com.",
	"sharethis.com.",
	"sharethrough.com.",
	"shiverscissors.com.",
	"sift.com.",
	"siftscience.com.",
	"singular.net.",
	"smartadserver.com.",
	"smilewanted.com.",
	"ssacdn.com.",
	"ssl-google-analytics.l.google.com.",
	"startmagazine.com.",
	"startscreen.com.",
	"stats.paypal.com.",
	"stickyadstv.com.",
	"supersonicads.com.",
	"supplyframe.com.",
	"szbnnqyqn.com.",
	"taboola.com.",
	"taboolanews.com.",
	"tagdeliver.com.",
	"telemetry.algolia.com.",
	"thestartmagazine.com.",
	"tifacciounregalo.com.",
	"typekit.net.",
	"unity3d.com.",
	"unseenreport.com.",
	"vungle.com.",
	"waggonerfoulpillow.com.",
	"wbbjrczne.com.",
	"wcsvjqcrpqqft.com.",
	"webengage.com.",
	"webtrekk.net.",
	"wt-safetag.com.",
	"www-googletagmanager.l.google.com.",
	"yieldmo.com.",
	"yourwebbars.com.",
	"zemanta.com.",
	"zorosrv.com.",
	"zqtk.net.",
}
