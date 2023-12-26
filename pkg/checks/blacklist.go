package checks

import (
	"strings"
)

func CheckDomain(domain string) bool {
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
	"corp.",
	"domain.",
	"example.",
	"home.",
	"host.",
	"internal.",
	"intranet.",
	"invalid.",
	"lan.",
	"local.",
	"localdomain.",
	"localhost.",
	"onion.",
	"private.",
	"rete.loc.",
	"test.",
	// China TLD
	"cn.",
	"xn--fiqs8s.",
	"xn--fiqz9s.",
	"hk.",
	"xn--j6w193g.",
	"mo.",
	"xn--mix082f.",
	"xn--mix891f.",
	"cyou.",
	// China domains
	"126.net.",
	"163.com.",
	"360buyimg.com.",
	"58.com.",
	"ali213.net.",
	"aliapp.org.",
	"alibabagroup.com.",
	"alicdn.com.",
	"aliespress.us.",
	"aliexpress.com.",
	"aliexpress.us.",
	"alimama.com.",
	"alipay.com.",
	"alipayobjects.com.",
	"alitrip.com.",
	"alixixi.com.",
	"aliyun.com.",
	"aliyuncs.com.",
	"amap.com.",
	"amemv.com.",
	"baidu.com.",
	"bcebos.com.",
	"bdimg.com.",
	"bdstatic.com.",
	"biliapi.net.",
	"bilibili.com.",
	"bilivideo.com.",
	"boom.skin.",
	"byted.org.",
	"bytedance.com.",
	"bytefcdn-oversea.com.",
	"byteoversea.com.",
	"china.com.",
	"csdn.net.",
	"ctrip.com.",
	"dbank.com.",
	"dbankcdn.com.",
	"dbankcloud.com.",
	"dbankcloud.eu.",
	"dbankedge.net.",
	"dianping.com.",
	"didispace.com.",
	"doh.pub.",
	"douyincdn.com.",
	"douyinliving.com.",
	"douyinpic.com.",
	"douyinvod.com.",
	"ecombdapi.com.",
	"etoote.com.",
	"freedidi.com.",
	"gifshow.com.",
	"harmonyos.com.",
	"hdslb.com.",
	"heytapmobi.com.",
	"hicloud.com.",
	"hsrtd.club.",
	"huanqiu.com.",
	"huawei.asia.",
	"huawei.com.",
	"huaweistatic.com.",
	"hwccpc.com.",
	"ibytedtos.com.",
	"ibyteimg.com.",
	"ifeng.com.",
	"inkuai.com.",
	"iqiyi.com.",
	"ixigua.com.",
	"jd.com.",
	"kekeshici.com.",
	"kuiniuca.com.",
	"kwaicdn.com.",
	"kwimgs.com.",
	"live-voip.com.",
	"meituan.com.",
	"meituan.net.",
	"mgtv.com.",
	"mi-img.com.",
	"mi.com.",
	"miui.com.",
	"mmstat.com.",
	"muscdn.com.",
	"musical.ly.",
	"musically.ly.",
	"myhuaweicloud.com.",
	"ndcpp.com.",
	"nuomi.com.",
	"onethingpcs.com.",
	"opera.com.",
	"pinduoduo.com.",
	"pstatp.com.",
	"qingting.fm.",
	"qplus.com.",
	"qq.com.",
	"qunar.com.",
	"qyer.com.",
	"sandai.net.",
	"servicewechat.com.",
	"snssdk.com.",
	"sogou.com.",
	"sohu.com.",
	"taobao.com.",
	"tencent-cloud.net.",
	"tencent.com.",
	"the-best-airport.com.",
	"tieba.com.",
	"tiktok.com.",
	"tiktok.net.",
	"tiktok.org.",
	"tiktokcdn.com.",
	"tiktokv.com.",
	"tmall.com.",
	"umeng.co.",
	"umeng.com.",
	"umengcloud.com.",
	"video-voip.com.",
	"vmall.com.",
	"wechat.com.",
	"wechat.org.",
	"wechatapp.com.",
	"weibo.com.",
	"xiachufang.com.",
	"xiami.com.",
	"xiaohongshu.com.",
	"xiaomi.com.",
	"xiaomi.net.",
	"xiaomiyoupin.com.",
	"ximalaya.com.",
	"xinhuanet.com.",
	"xunlei.com.",
	"xycdn.com.",
	"yinyuetai.com.",
	"youku.com.",
	"yximgs.com.",
	"zhihu.com.",
	"zijieapi.com.",
	// Russia TLD
	"ru.",
	"su.",
	"xn--p1ai.",
	"by.",
	"xn--90ais.",
	"yandex.",
	// Russia domains
	"2miners.com.",
	"adhigh.net.",
	"adlook.me.",
	"avito.st.",
	"betweendigital.com.",
	"bumlam.com.",
	"buzzoola.com.",
	"evraz.com.",
	"hybrid.ai.",
	"kas-labs.com.",
	"kaspersky-labs.com.",
	"kaspersky.com.",
	"kaspersky.it.",
	"labkas.com.",
	"metalloinvest.com.",
	"moe.video.",
	"mradx.net.",
	"mycdn.me.",
	"n1mk.com.",
	"nangs.org.",
	"nornickel.com.",
	"otm-r.com.",
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
	"yandex.com.",
	"yandex.eu.",
	"yandex.net.",
	"yastatic.net.",
	// Iran TLD
	"ir.",
	"xn--mgba3a4f16a.",
	// Iran domains
	"agah.com.",
	"aparat.com.",
	"bale.ai.",
	"cheryl.lol.",
	"digikala.com.",
	"digimovie.vip.",
	"digimoviez.com.",
	"eitaa.com.",
	"emofid.com.",
	"filimo.com.",
	"irene.lat.",
	"jocelyn.beauty.",
	"kittiy.sbs.",
	"lilliane.autos.",
	"maggiy.shop.",
	"marary.skin.",
	"maxin.beauty.",
	"mediaad.org.",
	"najva.com.",
	"nicoler.lol.",
	"nina.bond.",
	"nolra.cyou.",
	"qbtemmjn.site.",
	"rahavard365.com.",
	"sanjagh.com.",
	"sunista.info.",
	"telewebion.com.",
	"torob.com.",
	"tsetmc.com.",
	"varzesh3.com.",
	"xrrednbs.site.",
	"yektanet.com.",
	// North Korea
	"kp.",
	// Blocks
	"12ezo5v60.com.",
	"1drv.com.",
	"1ecosolution.it.",
	"1phads.com.",
	"1rx.io.",
	"24network.it.",
	"2mdn.net.",
	"2o7.net.",
	"33across.com.",
	"360yield.com.",
	"3lift.com.",
	"4chan-ads.org.",
	"4dsply.com.",
	"6bgaput9ullc.com.",
	"6g0blqi1541polz4n0kjvwo1kjl5tcx30.xyz.",
	"aadrm.com.",
	"aatkit.com.",
	"abtasty.com.",
	"acompli.net.",
	"activate.com.",
	"ad-center.com.",
	"ad-delivery.net.",
	"ad-stir.com.",
	"ad-sys.com.",
	"ad120m.com.",
	"ad127m.com.",
	"ad131m.com.",
	"ad132m.com.",
	"ad4mat.it.",
	"adadvisor.net.",
	"adasta.it.",
	"adcash.com.",
	"adcolony.com.",
	"addresseepaper.com.",
	"addthis.com.",
	"addthiscdn.com.",
	"addthisedge.com.",
	"addtoany.com.",
	"adform.com.",
	"adform.net.",
	"adformdsp.net.",
	"adj.st.",
	"adjust.com.",
	"adjust.net.in.",
	"adjust.world.",
	"adkaora.space.",
	"adlightning.com.",
	"admantx.com.",
	"admob.com.",
	"adnetworkperformance.com.",
	"adnext.it.",
	"adnexus.net.",
	"adnxs-simple.com.",
	"adnxs.com.",
	"adobe.com.",
	"adobe.it.",
	"adobedtm.com.",
	"adobeereg.com.",
	"adpartner.it.",
	"adplay.it.",
	"adpointrtb.com.",
	"adresponse.it.",
	"ads-twitter.com.",
	"ads.linkedin.com.",
	"adsafeprotected.com.",
	"adsco.re.",
	"adsensecustomsearchads.com.",
	"adservice.google.com.",
	"adsmoloco.com.",
	"adspro.it.",
	"adsrvr.org.",
	"adswizz.com.",
	"adtech.de.",
	"adtelligent.com.",
	"adtng.com.",
	"advertising.com.",
	"advgo.net.",
	"adzcore.com.",
	"aka.ms.",
	"amazon-adsystem.com.",
	"ammadv.it.",
	"amplitude.com.",
	"ampproject.net.",
	"ampproject.org.",
	"analytics.edgekey.net.",
	"analytics.githubassets.com.",
	"angsrvr.com.",
	"aniview.com.",
	"annunciadv.com.",
	"anydesk.com.",
	"app-measurement.com.",
	"app-measurement.com/sdk-exp.",
	"appboy.eu.",
	"appcenter.ms.",
	"applovin.com.",
	"applvn.com.",
	"appsflyer.com.",
	"appsflyersdk.com.",
	"artoukfarepu.org.",
	"asbgdfxrau.com.",
	"asmetotreatwab.com.",
	"assets-yammer.com.",
	"aswpsdkeu.com.",
	"atdmt.com.",
	"ayads.co.",
	"azurerms.com.",
	"banquetunarmedgrater.com.",
	"barscreative1.com.",
	"bbrdbr.com.",
	"bing.com.",
	"bing.net.",
	"bizible.com.",
	"bizzabo.com.",
	"blueconic.com.",
	"blueconic.net.",
	"bluekai.com.",
	"bodis.com.",
	"bounceexchange.com.",
	"branch.io.",
	"brandmetrics.com.",
	"braze-images.com.",
	"braze.com.",
	"braze.eu.",
	"brightcove.com.",
	"brightcove.net.",
	"browser-intake-datadoghq.com.",
	"brznetwork.com.",
	"bshrdr.com.",
	"btloader.com.",
	"bugsnag.com.",
	"buqkrzbrucz.com.",
	"camminachetipassa.it.",
	"casalemedia.com.",
	"cationinin.com.",
	"cbjqbkacuxjw.com.",
	"cdn1ve3zg.com.",
	"cdn4ads.com.",
	"cedexis-radar.net.",
	"cedexis.com.",
	"chartbeat.com.",
	"chartbeat.net.",
	"chatwoot.com.",
	"circulationrefill.com.",
	"clickcount.pw.",
	"clickiocmp.com.",
	"clickpoint.it.",
	"clickwave.media.",
	"closed.services.",
	"cloudappsecurity.com.",
	"cloudflareinsights.com.",
	"collector.github.com.",
	"collisionimpulsivejumpy.com.",
	"commento.io.",
	"commentsmodule.com.",
	"comscore.com.",
	"comscoreresearch.com.",
	"concert.io.",
	"conde.digital.",
	"confiant-integrations.global.ssl.fastly.net.",
	"confiant-integrations.net.",
	"connatix.com.",
	"connexity.com.",
	"connexity.net.",
	"consensu.org.",
	"consigli.it.",
	"convertexperiments.com.",
	"convertkit.com.",
	"conviva.com.",
	"cookiebot.com.",
	"cookiehub.net.",
	"cookielaw.org.",
	"coolbearsdaily54.com.",
	"cortana.ai.",
	"coveo.com.",
	"crashlogs.whatsapp.net.",
	"crashlytics.com.",
	"crashlyticsreports-pa.googleapis.com.",
	"crazyegg.com.",
	"creative-bars1.com.",
	"creative-serving.com.",
	"creativecdn.com.",
	"criteo.com.",
	"criteo.net.",
	"crittercism.com.",
	"crowdsignal.com.",
	"cruel-national.pro.",
	"crwdcntrl.net.",
	"cxense.com.",
	"czyasezpvs.com.",
	"d365ccafpi.com.",
	"datadoghq-browser-agent.com.",
	"declareave.com.",
	"definedlaunching.com.",
	"demdex.net.",
	"detectportal.firefox.com.",
	"directservices.it.",
	"discoveryfeed.org.",
	"displayvertising.com.",
	"disqus.com.",
	"disquscdn.com.",
	"districtm.io.",
	"dmca.com.",
	"dns.google.",
	"dns.google.com.",
	"dotmetrics.net.",
	"doubleclick.net.",
	"driftt.com.",
	"dynamics.com.",
	"dynamicyield.com.",
	"e-planning.net.",
	"ediscom.it.",
	"emukentsiwo.org.",
	"enslavequalities.com.",
	"eslzcjnlkepoow.com.",
	"etyequiremu.org.",
	"eventuallypropagandametal.com.",
	"evidon.com.",
	"evolutionadv.it.",
	"example.org.",
	"exelator.com.",
	"exmarketplace.com.",
	"exp-tas.com.",
	"ezoic.com.",
	"ezoic.net.",
	"facebook.com.",
	"facebook.net.",
	"fastly-insights.com.",
	"fbcdn.net.",
	"featuregates.org.",
	"fedfceddu.com.",
	"firebaselogging-pa.googleapis.com.",
	"flowshaft.com.",
	"fontawesome.com.",
	"fonts.googleapis.com.",
	"fonts.gstatic.com.",
	"fonts.shopifycdn.com.",
	"forfeitsubscribe.com.",
	"forlumineoner.com.",
	"fqdwrgbbkmlbh.com.",
	"fraud0.com.",
	"friendshipmale.com.",
	"ftjcfx.com.",
	"fundingchoicesmessages.google.com.",
	"fvcwqkkqmuv.com.",
	"gaconnector.com.",
	"gameloft.com.",
	"gatekeeperconsent.com.",
	"geoedge.be.",
	"getadmiral.com.",
	"getkoala.com.",
	"getpocket.com.",
	"ghabovethec.info.",
	"gigya.com.",
	"gl-product-analytics.com.",
	"glyph.medium.com.",
	"gmxvmvptfm.com.",
	"google-analytics.com.",
	"googleadservices.com.",
	"googlesyndication.com.",
	"googletagmanager.com.",
	"googletagservices.com.",
	"grow.me.",
	"gumgum.com.",
	"hebfcxdchubdbs.com.",
	"histats.com.",
	"hockeyapp.net.",
	"honorablehalt.com.",
	"hotjar.com.",
	"hs-analytics.net.",
	"hsadspixel.net.",
	"id5-sync.com.",
	"iinnbewhna.com.",
	"imasdk.googleapis.com.",
	"img-taboola.com.",
	"imitrkn.net.",
	"improving.duckduckgo.com.",
	"imrworldwide.com.",
	"inlandexaminerinterrogate.com.",
	"inmobi.com.",
	"inner-active.mobi.",
	"instagram.com.",
	"instructorloneliness.com.",
	"insurads.com.",
	"intake-analytics.wikimedia.org.",
	"intellipopup.com.",
	"intelliscapesolutions.com.",
	"intercom.io.",
	"intercomcdn.com.",
	"involvingaged.com.",
	"ioladv.it.",
	"iperbanner.com.",
	"ipv4only.arpa.",
	"isgprivacy.cbsi.com.",
	"iubenda.com.",
	"jads.co.",
	"jf8d0dlc.com.",
	"jhldtogerghottulering.info.",
	"jifo.co.",
	"jpyxqeysh.com.",
	"juiceadv.com.",
	"juicyads.com.",
	"jumptap.com.",
	"ketch.com.",
	"ketchjs.com.",
	"kijiji.it.",
	"koala.live.",
	"ku2d3a7pa8mdi.com.",
	"ku42hjr2e.com.",
	"kustomerapp.com.",
	"lasubqueries.com.",
	"laughedaffront.com.",
	"launchdarkly.com.",
	"leaderaffiliation.com.",
	"leonardoadv.it.",
	"limurol.com.",
	"liteanalytics.com.",
	"live.com.",
	"live.net.",
	"logrocket-assets.io.",
	"ltmsphrcl.net.",
	"lukecomparetwo.com.",
	"lunalabs.io.",
	"lwonclbench.com.",
	"lync.com.",
	"main-instinct.com.",
	"mambasms.com.",
	"marketingcloudapis.com.",
	"marketo.net.",
	"matterlytics.com.",
	"mediavine.com.",
	"memoinsights.com.",
	"merequartz.com.",
	"mgsn.it.",
	"microsoft.com.",
	"microsoftonline-p.com.",
	"microsoftonline-p.net.",
	"microsoftonline.com.",
	"microsoftstream.com.",
	"mlepvbgowvzt.com.",
	"mnaspm.com.",
	"moatads.com.",
	"mobileapptracking.com.",
	"mobilesecuremail.com.",
	"moengage.com.",
	"mparticle.com.",
	"mpi-internal.com.",
	"mpianalytics.com.",
	"mrf.io.",
	"msads.net.",
	"msauth.net.",
	"msauthimages.net.",
	"msft.net.",
	"msftauth.net.",
	"msftauthimages.net.",
	"msftidentity.com.",
	"msidentity.com.",
	"msocdn.com.",
	"mstea.ms.",
	"mycookies.it.",
	"mysura.it.",
	"myvisualiq.net.",
	"nativery.com.",
	"neodatagroup.com.",
	"neverbounce.com.",
	"newrelic.com.",
	"ninanceenab.com.",
	"nocookie.net.",
	"nondescriptnote.com.",
	"novadv.com.",
	"npttech.com.",
	"nr-data.net.",
	"o365weve.com.",
	"oaspapps.com.",
	"office.com.",
	"office.net.",
	"office365.com.",
	"ojuhjcmhemvs.com.",
	"okta.com.",
	"omnitagjs.com.",
	"omtrdc.net.",
	"one.one.",
	"onedrive.com.",
	"onenote.com.",
	"onenote.net.",
	"onestore.ms.",
	"onetag-sys.com.",
	"onetrust.com.",
	"onetrust.io.",
	"online-metrix.net.",
	"onmicrosoft.com.",
	"opecloud.com.",
	"openx.net.",
	"opsource.net.",
	"optimizely.com.",
	"orgotitedu.info.",
	"osano.com.",
	"outbrain.com.",
	"outbrainimg.com.",
	"outlook.com.",
	"outlookmobile.com.",
	"owox.com.",
	"padfungusunless.com.",
	"parse.ly.",
	"parsely.com.",
	"partypartners.it.",
	"passport.net.",
	"payclick.it.",
	"pemsrv.com.",
	"perfectmarket.com.",
	"permutive.app.",
	"permutive.com.",
	"pghub.io.",
	"phonefactor.net.",
	"piano.io.",
	"piwik.pro.",
	"pixel.wp.com.",
	"platform.linkedin.com.",
	"platform.twitter.com.",
	"plausible.io.",
	"play.google.com.",
	"playfabapi.com.",
	"pogothere.xyz.",
	"polldaddy.com.",
	"powerapps.com.",
	"powerplatform.com.",
	"prebid.org.",
	"privacy-center.org.",
	"privacy-mgmt.com.",
	"privacymanager.io.",
	"proftrafficcounter.com.",
	"proxy-safebrowsing.googleapis.com.",
	"pub.network.",
	"pubmatic.com.",
	"pubmnet.com.",
	"pubnative.net.",
	"pubtech.ai.",
	"pushio.com.",
	"quantcast.com.",
	"quantcount.com.",
	"quantserve.com.",
	"raddoppia-bitcoin.click.",
	"ravelin.click.",
	"rcsmetrics.it.",
	"realsrv.com.",
	"repentantsympathy.com.",
	"researchnow.com.",
	"revolt.chat.",
	"richaudience.com.",
	"rqmob.com.",
	"rtbix.xyz.",
	"rtbuzz.net.",
	"rtmark.net.",
	"rubiconproject.com.",
	"run-syndicate.com.",
	"runative-syndicate.com.",
	"sail-horizon.com.",
	"sascdn.com.",
	"sb89347.com.",
	"sbanner.com.",
	"sc-cdn.net.",
	"sc-gw.com.",
	"sc-jpl.com.",
	"sc-prod.net.",
	"scambiobanner.it.",
	"scambiobanner.org.",
	"scambiobanner.tv.",
	"scambiositi.com.",
	"scaredsnakes.com.",
	"scarletwood.com.",
	"schip.io.",
	"scorecardresearch.com.",
	"sda.fyi.",
	"seadform.net.",
	"seedtag.com.",
	"sembox.it.",
	"sendbird.com.",
	"sentry-cdn.com.",
	"sentry.io.",
	"seolabadv.it.",
	"servedby-buysellads.com.",
	"servenobid.com.",
	"serving-sys.com.",
	"sfbassets.com.",
	"sfx.ms.",
	"sharepointonline.com.",
	"sharethis.com.",
	"sharethrough.com.",
	"shiverscissors.com.",
	"sift.com.",
	"siftscience.com.",
	"singular.net.",
	"skype.com.",
	"skypeassets.com.",
	"skypeforbusiness.com.",
	"smartadserver.com.",
	"smilewanted.com.",
	"snap.com.",
	"snapads.com.",
	"snapchat.com.",
	"snapkit.com.",
	"sneerattendingconverted.com.",
	"sophi.io.",
	"spacetraff.com.",
	"spmailtechnol.com.",
	"spot.im.",
	"spotim.martket.",
	"srvs.site.",
	"ssacdn.com.",
	"ssl-google-analytics.l.google.com.",
	"staffboozerenamed.com.",
	"staffhub.ms.",
	"startmagazine.com.",
	"startscreen.com.",
	"stats.paypal.com.",
	"stats.wp.com.",
	"statsig.com.",
	"statsigapi.net.",
	"statsigapi.net.",
	"statsigcdn.com.",
	"stickyadstv.com.",
	"stripcdn.com.",
	"strpst.com.",
	"supersonicads.com.",
	"supplyframe.com.",
	"sway-cdn.com.",
	"sway-extensions.com.",
	"sway.com.",
	"swonqjzbc.com.",
	"szbnnqyqn.com.",
	"taboola.com.",
	"taboolanews.com.",
	"tagdeliver.com.",
	"talkscreativity.com.",
	"teads.tv.",
	"tealium.com.",
	"telemetry.algolia.com.",
	"the-ozone-project.com.",
	"thestartmagazine.com.",
	"tifacciounregalo.com.",
	"tinypass.com.",
	"tiqcdn.com.",
	"tncid.app.",
	"trackwilltrk.com.",
	"trendemon.com.",
	"triplelift.com.",
	"troopsassistedstupidity.com.",
	"tru.am.",
	"trustarc.com.",
	"tsyndicate.com.",
	"tvpixel.com.",
	"typekit.net.",
	"unity3d.com.",
	"unseenreport.com.",
	"usercentrics.com.",
	"usercentrics.eu.",
	"userzoom.com.",
	"vhkbvpbuhwon.com.",
	"virtualearth.net.",
	"visualwebsiteoptimizer.com.",
	"vpn-content.net.",
	"vpn-world.com.",
	"vungle.com.",
	"waggonerfoulpillow.com.",
	"washealinginduced.com.",
	"wbbjrczne.com.",
	"wcsvjqcrpqqft.com.",
	"weaktongue.com.",
	"webcontentassessor.com.",
	"webengage.com.",
	"webtrekk.net.",
	"widgets.wp.com.",
	"windows-ppe.net.",
	"windows.com.",
	"windows.net.",
	"wt-safetag.com.",
	"www-googletagmanager.l.google.com.",
	"ximxim.com.",
	"xlivrdr.com.",
	"yahoosandbox.com.",
	"yammer.com.",
	"yammerusercontent.com.",
	"ybb-network.com.",
	"ybbserver.com.",
	"yieldmo.com.",
	"yobee.it.",
	"yourwebbars.com.",
	"zemanta.com.",
	"zoom.us.",
	"zorosrv.com.",
	"zqtk.net.",
}
