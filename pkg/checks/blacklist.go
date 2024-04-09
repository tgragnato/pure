package checks

var blacklist = []string{
	// Reserved
	"arpa.",
	"corp.",
	"domain.",
	"example.",
	"farm.",
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
	"bdurl.net.",
	"biliapi.com.",
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
	"dnspod.com.",
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
	"ksapisrv.com.",
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
	"mjt000.com.",
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
	"tencent-cloud.com.",
	"tencent-cloud.net.",
	"tencent.com.",
	"tencentcloud.com.",
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
	"my.com.",
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
	"1ts19.top.",
	"20000.click.",
	"24network.it.",
	"26485.top.",
	"2mdn.net.",
	"2o7.net.",
	"33across.com.",
	"360yield.com.",
	"3lift.com.",
	"4chan-ads.org.",
	"4dex.io.",
	"4dsply.com.",
	"4strokemedia.com.",
	"5i68sbhin.com.",
	"6bgaput9ullc.com.",
	"6g0blqi1541polz4n0kjvwo1kjl5tcx30.xyz.",
	"6sc.co.",
	"9l3s3fnhl.com.",
	"a-mo.net.",
	"a-mx.com.",
	"a2z.com.",
	"aadrm.com.",
	"aatkit.com.",
	"abchygmsaftnrr.xyz.",
	"abinsula.com.",
	"ablecolony.com.",
	"abletoprese.org.",
	"abtasty.com.",
	"acertb.com.",
	"acompli.net.",
	"activate.com.",
	"acuityplatform.com.",
	"ad-center.com.",
	"ad-delivery.net.",
	"ad-m.asia.",
	"ad-stir.com.",
	"ad-sys.com.",
	"ad.gt.",
	"ad120m.com.",
	"ad127m.com.",
	"ad131m.com.",
	"ad132m.com.",
	"ad4mat.it.",
	"adadvisor.net.",
	"adapex.io.",
	"adasta.it.",
	"adcash.com.",
	"adcolony.com.",
	"addlnk.com.",
	"addresseepaper.com.",
	"addthis.com.",
	"addthiscdn.com.",
	"addthisedge.com.",
	"addtoany.com.",
	"adentifi.com.",
	"adforcast.com.",
	"adform.com.",
	"adform.net.",
	"adformdsp.net.",
	"adglare.net.",
	"adgrid.io.",
	"adgrx.com.",
	"adhaven.com.",
	"adition.com.",
	"adj.st.",
	"adjust.com.",
	"adjust.net.in.",
	"adjust.world.",
	"adkaora.space.",
	"adkernel.com.",
	"adlightning.com.",
	"admantx.com.",
	"admedo.com.",
	"admixer.net.",
	"admob.com.",
	"adnami.io.",
	"adnetworkperformance.com.",
	"adnext.it.",
	"adnexus.net.",
	"adnxs-simple.com.",
	"adnxs.com.",
	"adobe.com.",
	"adobe.it.",
	"adobedtm.com.",
	"adobeereg.com.",
	"adotmob.com.",
	"adpartner.it.",
	"adpartner.pro.",
	"adplay.it.",
	"adpointrtb.com.",
	"adpushup.com.",
	"adresponse.it.",
	"adroll.com.",
	"ads-api.twitter.com.",
	"ads-twitter.com.",
	"ads.linkedin.com.",
	"ads.youtube.com.",
	"adsafeprotected.com.",
	"adsafety.net.",
	"adscale.de.",
	"adsco.re.",
	"adsensecustomsearchads.com.",
	"adservice.google.ca.",
	"adservice.google.com.",
	"adservice.google.it.",
	"adskeeper.com.",
	"adsmoloco.com.",
	"adspro.it.",
	"adspsp.com.",
	"adsrvr.org.",
	"adstanding.com.",
	"adswizz.com.",
	"adsymptotic.com.",
	"adtdp.com.",
	"adtech.de.",
	"adtelligent.com.",
	"adthrive.com.",
	"adtng.com.",
	"adultfriendfinder.com.",
	"adv.ilpost.it.",
	"adv.ilriformista.it.",
	"advertising.amazon.it.",
	"advertising.com.",
	"advgo.net.",
	"adxbid.info.",
	"adxpremium.services.",
	"adzcore.com.",
	"affec.tv.",
	"afilliateapp.com.",
	"afodreet.net.",
	"agamaevascla.top.",
	"agkn.com.",
	"aidata.io.",
	"airbrake.io.",
	"aka.ms.",
	"akstat.io.",
	"alas4kanmfa6a4mubte.com.",
	"alcmpn.com.",
	"alexametrics.com.",
	"amazon-adsystem.com.",
	"amgtui.com.",
	"amirlabd.com.",
	"ammadv.it.",
	"ampaeservices.com.",
	"amplitude.com.",
	"ampproject.net.",
	"ampproject.org.",
	"amung.us.",
	"analysis.fi.",
	"analytics-tracking.meetup.com.",
	"analytics.apache.org.",
	"analytics.astalegale.net.",
	"analytics.edgekey.net.",
	"analytics.githubassets.com.",
	"analytics.google.com.",
	"analytics.google.it.",
	"analytics.kaltura.com.",
	"analytics.pointdrive.linkedin.com.",
	"angsrvr.com.",
	"aniview.com.",
	"annunciadv.com.",
	"anydesk.com.",
	"app-measurement.com.",
	"app-measurement.com/sdk-exp.",
	"app.link.",
	"appboy.eu.",
	"appboycdn.com.",
	"appcenter.ms.",
	"appconsent.io.",
	"appcues.com.",
	"appier.net.",
	"applovin.com.",
	"applvn.com.",
	"appsflyer.com.",
	"appsflyersdk.com.",
	"arqsafhutlam.com.",
	"artoukfarepu.org.",
	"arvigorothan.com.",
	"asbgdfxrau.com.",
	"ascendex.com.",
	"asdjuwobnagm.com.",
	"asmetotreatwab.com.",
	"assemblyexchange.com.",
	"assets-yammer.com.",
	"assoc-amazon.com.",
	"aswpsdkeu.com.",
	"aswpsdkus.com.",
	"atdmt.com.",
	"aterhouseoyop.com.",
	"atmtd.com.",
	"attn.tv.",
	"authenticsystemtraffic.com.",
	"automatad.com.",
	"automizely.com.",
	"avo.app.",
	"avocet.io.",
	"ayads.co.",
	"azurerms.com.",
	"backuprabbit.com.",
	"badgerstat.com.",
	"bande2az.com.",
	"banquetunarmedgrater.com.",
	"barscreative1.com.",
	"basis.net.",
	"bazaarvoice.com.",
	"bbrdbr.com.",
	"beacons.ai.",
	"bestoftheat.com.",
	"bet365.it.",
	"betotodilea.com.",
	"bfmio.com.",
	"bidmatic.io.",
	"bidr.io.",
	"bidswitch.net.",
	"binaryborrowedorganized.com.",
	"bing.com.",
	"bing.net.",
	"bizible.com.",
	"bizographics.com.",
	"bizrate.com.",
	"bizzabo.com.",
	"bkrtx.com.",
	"blismedia.com.",
	"blockadsnot.com.",
	"bluecava.com.",
	"blueconic.com.",
	"blueconic.net.",
	"bluekai.com.",
	"bmaoghgwyhi.com.",
	"bobabillydirect.org.",
	"bodis.com.",
	"bounceexchange.com.",
	"brainlyads.com.",
	"branch.io.",
	"brandmetrics.com.",
	"braze-images.com.",
	"braze.com.",
	"braze.eu.",
	"brealtime.com.",
	"brightcove.com.",
	"brightcove.net.",
	"bronto.com.",
	"browser-intake-datadoghq.com.",
	"browser-update.org.",
	"browsiprod.com.",
	"brunchcreatesenses.com.",
	"brznetwork.com.",
	"bshrdr.com.",
	"btloader.com.",
	"btstatic.com.",
	"bttrack.com.",
	"btttag.com.",
	"bugsnag.com.",
	"buqkrzbrucz.com.",
	"buysellads.com.",
	"bvpxtrck.com.",
	"bygliscortor.com.",
	"bytogeticr.com.",
	"caicaipoop.com.",
	"caizutoh.xyz.",
	"cam4tracking.com.",
	"cameesse.net.",
	"camminachetipassa.it.",
	"capaciousdrewreligion.com.",
	"capletstyldia.com.",
	"casalemedia.com.",
	"casino-zilla.com.",
	"cationinin.com.",
	"cbjqbkacuxjw.com.",
	"cdn1ve3zg.com.",
	"cdn4ads.com.",
	"cdnfonts.com.",
	"cdnwidget.com.",
	"cedexis-radar.akamaized.net.",
	"cedexis-radar.net.",
	"cedexis-test.akamaized.net.",
	"cedexis-test.com.",
	"cedexis.com.",
	"chartbeat.com.",
	"chartbeat.net.",
	"chatwoot.com.",
	"churnedflames.top.",
	"ciaopeople.it.",
	"circulationrefill.com.",
	"civiccomputing.com.",
	"clarity.ms.",
	"classicguarantee.pro.",
	"clearbit.com.",
	"clearbitjs.com.",
	"clearbitscripts.com.",
	"cleverwebserver.com.",
	"clickagy.com.",
	"clickcertain.com.",
	"clickcount.pw.",
	"clickiocdn.com.",
	"clickiocmp.com.",
	"clickpoint.it.",
	"clickwave.media.",
	"closed.services.",
	"cloudappsecurity.com.",
	"cloudflareinsights.com.",
	"cloudfrale.com.",
	"cloudfront.net.",
	"cloudsponcer.com.",
	"clrstm.com.",
	"cmpct.info.",
	"cohesionapps.com.",
	"colarak.com.",
	"collector.github.com.",
	"collisionimpulsivejumpy.com.",
	"colossusssp.com.",
	"commander1.",
	"commento.io.",
	"commentsmodule.com.",
	"company-target.com.",
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
	"consentframework.com.",
	"consentmanager.net.",
	"consigli.it.",
	"content-storage-download.googleapis.com.",
	"content-storage-upload.googleapis.com.",
	"contentabc.com.",
	"contentsquare.net",
	"contextweb.com.",
	"convertexperiments.com.",
	"convertkit.com.",
	"convertlanguage.com.",
	"conviva.com.",
	"cookie-script.com.",
	"cookiebot.com.",
	"cookiebot.eu.",
	"cookiefirst.com.",
	"cookiehub.eu.",
	"cookiehub.net.",
	"cookieinformation.com.",
	"cookielaw.org.",
	"cookiepro.com.",
	"cookieyes.com.",
	"cooladata.com.",
	"coolbearsdaily54.com.",
	"coollabs.io.",
	"coosync.com.",
	"cootlogix.com.",
	"cortana.ai.",
	"coveo.com.",
	"cpx.to.",
	"crashlogs.whatsapp.net.",
	"crashlytics.com.",
	"crashlyticsreports-pa.googleapis.com.",
	"crazyegg.com.",
	"creative-bars1.com.",
	"creative-serving.com.",
	"creativecdn.com.",
	"crisp.chat.",
	"criteo.com.",
	"criteo.net.",
	"crittercism.com.",
	"crowdsignal.com.",
	"crsspxl.com.",
	"cruel-national.pro.",
	"crwdcntrl.net.",
	"cryingforanythi.com.",
	"ctnsnet.com.",
	"ctrtrk.com.",
	"customads.co.",
	"cxense.com.",
	"czyasezpvs.com.",
	"d365ccafpi.com.",
	"dadalytics.it.",
	"daicagrithi.com.",
	"datadoghq-browser-agent.com.",
	"datadome.co.",
	"datatechone.com.",
	"dcyqtufkopp.com.",
	"declareave.com.",
	"deepintent.com.",
	"definedlaunching.com.",
	"deliverimp.com.",
	"demandbase.com.",
	"demdex.net.",
	"descendentwringthou.com.",
	"detectportal.firefox.com.",
	"deywepri.com.",
	"directservices.it.",
	"discoveryfeed.org.",
	"displayvertising.com.",
	"disqus.com.",
	"disquscdn.com.",
	"disshipbikinis.com.",
	"districtm.io.",
	"dkbupulnfm.com.",
	"dmca.com.",
	"dmsktmld.com.",
	"dns.google.",
	"dns.google.com.",
	"domdex.com.",
	"domuipan.com.",
	"donnemagazine.it.",
	"dotmetrics.net.",
	"dotomi.com.",
	"doubleclick.net.",
	"doubleverify.com.",
	"downstairsnegotiatebarren.com.",
	"dpmsrv.com.",
	"dreamlab.pl.",
	"driftt.com.",
	"dsail-tech.com.",
	"dtpeqgfaps.com.",
	"dtscout.com.",
	"dukicationan.org.",
	"dungeonisosculptor.com.",
	"dvypar.com.",
	"dwin1.com.",
	"dynamics.com.",
	"dynamicyield.com.",
	"dyntrk.com.",
	"e-planning.net.",
	"e2ertt.com.",
	"echobox.com.",
	"ediscom.it.",
	"eergortu.net.",
	"ejixwhxmacu.xyz.",
	"ekbl.net.",
	"elasticad.net.",
	"eloqua.com.",
	"emukentsiwo.org.",
	"emxdgt.com.",
	"en25.com.",
	"engageplatform.com.",
	"ensighten.com.",
	"enslavequalities.com.",
	"eslzcjnlkepoow.com.",
	"ethicalads.io.",
	"etoro.com.",
	"etyequiremu.org.",
	"evecticvocoder.life.",
	"eveneraw.digital.",
	"eventuallypropagandametal.com.",
	"everesttech.net.",
	"evidon.com.",
	"evolutionadv.it.",
	"example.org.",
	"exdynsrv.com.",
	"exelator.com.",
	"exmarketplace.com.",
	"exoclick.com.",
	"exosrv.com.",
	"exp-tas.com.",
	"exponential.com.",
	"exposebox.com.",
	"extanalytics.com.",
	"extend.tv.",
	"eyeota.net.",
	"eyeviewads.com.",
	"ezodn.com.",
	"ezoic.com.",
	"ezoic.net.",
	"ezojs.com.",
	"ezzmvdyyzccx.com.",
	"facebook.com.",
	"facebook.net.",
	"fadingmummytuxedo.com.",
	"fallclk.com.",
	"fastclick.net.",
	"fastly-insights.com.",
	"fbcdn.net.",
	"fbsbx.com.",
	"fearlessfaucet.com.",
	"feathr.co.",
	"featuregates.org.",
	"fedfceddu.com.",
	"feedbackify.com.",
	"fg8dgt.com.",
	"firebaselogging-pa.googleapis.com.",
	"firstimpression.io.",
	"fkdhbmsss.com.",
	"flashtalking.com.",
	"flowshaft.com.",
	"fontawesome.com.",
	"fonts-api.wp.com.",
	"fonts.googleapis.com.",
	"fonts.gstatic.com.",
	"fonts.net.",
	"fonts.shopifycdn.com.",
	"fonts.wp.com.",
	"fontshare.com.",
	"foodblog.it.",
	"foresee.com.",
	"forfeitsubscribe.com.",
	"forlumineoner.com.",
	"forthemoonh.com.",
	"fqdwrgbbkmlbh.com.",
	"fqtag.com.",
	"fraud0.com.",
	"freshmarketer.com.",
	"freshmarketer.eu.",
	"freshpops.net.",
	"freshsales.io.",
	"friendshipmale.com.",
	"frydtyhhya.com.",
	"ftjcfx.com.",
	"fundingchoicesmessages.google.com.",
	"furthermoreimpetusscribble.com.",
	"fvcwqkkqmuv.com.",
	"fwmrm.net.",
	"fxkdv.com.",
	"gaconnector.com.",
	"gameloft.com.",
	"gaming-adult.com.",
	"gatekeeperconsent.com.",
	"gemius.pl.",
	"geoedge.be.",
	"geoip.businessinsider.com.",
	"geojs.io.",
	"geolocation.forbes.com.",
	"getadmiral.com.",
	"getkoala.com.",
	"getpocket.com.",
	"getscriptjs.com.",
	"getsitecontrol.com.",
	"gfx.ms.",
	"ghabovethec.info.",
	"gigya.com.",
	"gl-product-analytics.com.",
	"glyph.medium.com.",
	"gmdigital.it.",
	"gmxvmvptfm.com.",
	"go-mpulse.net.",
	"godpvqnszo.com.",
	"google-analytics.com.",
	"googleadservices.com.",
	"googleanalytics.com.",
	"googledatas.com.",
	"googleoptimize.com.",
	"googlesyndication.com.",
	"googletagmanager.com.",
	"googletagservices.com.",
	"gptkjrseu.com.",
	"gravitec.net.",
	"groorsoa.net.",
	"grow.me.",
	"growsumo.com.",
	"growthbook.io.",
	"grsm.io.",
	"gssprt.jp.",
	"gumgum.com.",
	"gwallet.com.",
	"happymuttere.org.",
	"heapanalytics.com.",
	"hebfcxdchubdbs.com.",
	"hellobar.com.",
	"heownouncillor.com.",
	"histats.com.",
	"hockeyapp.net.",
	"honorablehalt.com.",
	"hotjar.com.",
	"hprofits.com.",
	"href.li.",
	"hrrlyfdnxlzxe.com.",
	"hs-analytics.net.",
	"hs-banner.com.",
	"hs-scripts.com.",
	"hsadspixel.net.",
	"hsleadflows.net.",
	"hubapi.com.",
	"hubspot.com.",
	"hum.works.",
	"humiliatingregion.com.",
	"humpdecompose.com.",
	"hwpnocpctu.com.",
	"ib-ibi.com.",
	"id5-sync.com.",
	"idescargarapk.com.",
	"igodigital.com.",
	"iinnbewhna.com.",
	"iloptrex.com.",
	"imasdk.googleapis.com.",
	"img-taboola.com.",
	"imitrkn.net.",
	"impactradius-event.com.",
	"improving.duckduckgo.com.",
	"imrworldwide.com.",
	"indexww.com.",
	"inlandexaminerinterrogate.com.",
	"inmobi.com.",
	"inner-active.mobi.",
	"innovid.com.",
	"inoculateconsessionconsessioneuropean.com.",
	"insightexpressai.com.",
	"instagram.com.",
	"instana.io.",
	"instructorloneliness.com.",
	"insurads.com.",
	"intake-analytics.wikimedia.org.",
	"intellipopup.com.",
	"intelliscapesolutions.com.",
	"intentiq.com.",
	"intercom.io.",
	"intercomcdn.com.",
	"involvingaged.com.",
	"ioam.de.",
	"iol.it.",
	"ioladv.it.",
	"iolam.it.",
	"ip-api.com.",
	"iperbanner.com.",
	"ipify.org.",
	"ipinfo.io.",
	"ipredictive.com.",
	"ipv4only.arpa.",
	"iqtewqshjpk.xyz.",
	"iqzone.com.",
	"ironsrc.mobi.",
	"isgprivacy.cbsi.com.",
	"islerobserpent.com.",
	"ispot.tv.",
	"isprog.com.",
	"istoanaugrub.xyz.",
	"itsup.com.",
	"iubenda.com.",
	"ixiaa.com.",
	"jaavnacsdw.com.",
	"jads.co.",
	"jeekomih.com.",
	"jenivlcamyph.xyz.",
	"jf8d0dlc.com.",
	"jhldtogerghottulering.info.",
	"jifo.co.",
	"jkanpupquilcgjl.xyz.",
	"jlufbcef.com.",
	"jointag.com.",
	"jpyxqeysh.com.",
	"juiceadv.com.",
	"juicyads.com.",
	"jumptap.com.",
	"justpremium.com.",
	"jwpltx.com.",
	"jygotubvpyguak.com.",
	"k5a.io.",
	"kampyle.com.",
	"kargo.com.",
	"kataweb.it.",
	"kdprquajwnr.com.",
	"keapgypsite.website.",
	"ketadexchange.com.",
	"ketch.com.",
	"ketchjs.com.",
	"keywee.co.",
	"kijiji.it.",
	"klaviyo.com.",
	"koala.live.",
	"krxd.net.",
	"ku2d3a7pa8mdi.com.",
	"ku42hjr2e.com.",
	"kumtibsa.com.",
	"kustomerapp.com.",
	"labadena.com.",
	"ladsp.com.",
	"lasubqueries.com.",
	"laughedaffront.com.",
	"launchdarkly.com.",
	"lby2kd27c.com.",
	"lcwoewvvmhj.com.",
	"leaderaffiliation.com.",
	"leadingindication.pro.",
	"lemsoodol.com.",
	"leonardoadv.it.",
	"letopreseynatc.org.",
	"lfstmedia.com.",
	"liadm.com.",
	"liaoptse.net.",
	"lightboxcdn.com.",
	"lijit.com.",
	"likebtn.com.",
	"limurol.com.",
	"linksynergy.com.",
	"linkvertise.com.",
	"listrak.com.",
	"listrakbi.com.",
	"liteanalytics.com.",
	"live.com.",
	"live.net.",
	"livechatinc.com.",
	"lkqd.net.",
	"llinks.io.",
	"lockerdome.com.",
	"logrocket-assets.io.",
	"loopme.me.",
	"ltedafajhb.com.",
	"ltmsphrcl.net.",
	"luckyorange.com.",
	"lukecomparetwo.com.",
	"lunalabs.io.",
	"lwonclbench.com.",
	"lync.com.",
	"m0rsq075u.com.",
	"magellanotech.it.",
	"magsrv.com.",
	"main-instinct.com.",
	"mambasms.com.",
	"marinsm.com.",
	"markedoneofth.com.",
	"marketingcloudapis.com.",
	"marketo.com.",
	"marketo.net.",
	"mathtag.com.",
	"matomo.cloudfront.similarweb.io.",
	"matomo.similarweb.io.",
	"matterlytics.com.",
	"maxymiser.net.",
	"maze.co.",
	"mcizas.com.",
	"mcpuwpush.com.",
	"mdstats.info.",
	"mdundo.com.",
	"medallia.com.",
	"medallia.eu.",
	"media.net.",
	"media6degrees.com.",
	"mediaplex.com.",
	"mediarithmics.com.",
	"mediavine.com.",
	"mediawallahscript.com.",
	"melansida.com.",
	"memoinsights.com.",
	"merequartz.com.",
	"metastatus.com.",
	"metricswpsh.com.",
	"mfadsrvr.com.",
	"mgid.com.",
	"mgsn.it.",
	"micpn.com.",
	"microad.jp.",
	"microsoft.com.",
	"microsoftonline-p.com.",
	"microsoftonline-p.net.",
	"microsoftonline.com.",
	"microsoftstream.com.",
	"mileesidesu.org.",
	"miltlametta.com.",
	"minutemedia-prebid.com.",
	"mixpanel.com.",
	"mkkvprwskq.com.",
	"mktoresp.com.",
	"mkyyhfbpyd.xyz.",
	"ml314.com.",
	"mlepvbgowvzt.com.",
	"mnaspm.com.",
	"mniumlapsers.com.",
	"moatads.com.",
	"mobileapptracking.com.",
	"mobilesecuremail.com.",
	"moengage.com.",
	"monetate.net.",
	"monetizer.com.",
	"mookie1.com.",
	"mordoops.com.",
	"motorimagazine.it.",
	"mouseflow.com.",
	"mparticle.com.",
	"mpi-internal.com.",
	"mpianalytics.com.",
	"mrf.io.",
	"mrpdata.net.",
	"msads.net.",
	"msauth.net.",
	"msauthimages.net.",
	"mscoldness.com.",
	"msft.net.",
	"msftauth.net.",
	"msftauthimages.net.",
	"msftconnecttest.com.",
	"msftidentity.com.",
	"msidentity.com.",
	"msn.com.",
	"msocdn.com.",
	"mspznxahjjx.com.",
	"mspznxahjjx.com.",
	"mstea.ms.",
	"mutteredadis.org.",
	"mxpnl.com.",
	"mxptint.net.",
	"mycookies.it.",
	"myfreshworks.com.",
	"mysura.it.",
	"myunderthfe.info.",
	"myvisualiq.net.",
	"narrativ.com.",
	"native-track.com.",
	"nativery.com.",
	"navdmp.com.",
	"neintheworld.org.",
	"neodatagroup.com.",
	"neondata.cloud.",
	"nereserv.com.",
	"netmng.com.",
	"networkpccontrol.com.",
	"neverbounce.com.",
	"newrelic.com.",
	"ninanceenab.com.",
	"ninthdecimal.com.",
	"nocaudsomt.xyz.",
	"nocookie.net.",
	"noibu.com.",
	"nominalclck.name.",
	"nondescriptnote.com.",
	"notsy.io.",
	"novadv.com.",
	"npttech.com.",
	"nr-data.net.",
	"ntsiwoulukdli.org.",
	"ntv.io.",
	"nzsndovdjs.com.",
	"o365weve.com.",
	"oaspapps.com.",
	"obeahwidowed.digital.",
	"obeywish.com.",
	"ofdrapiona.com.",
	"office.com.",
	"office.net.",
	"office365.com.",
	"ojgabhavrxub.com.",
	"ojuhjcmhemvs.com.",
	"okta.com.",
	"olark.com.",
	"omappapi.com.",
	"omnitagjs.com.",
	"omtrdc.net.",
	"one.one.",
	"onedrive.com.",
	"onenote.com.",
	"onenote.net.",
	"onesignal.com.",
	"onestore.ms.",
	"onetag-sys.com.",
	"onetrust.com.",
	"onetrust.io.",
	"online-metrix.net.",
	"onmicrosoft.com.",
	"onthe.io.",
	"opecloud.com.",
	"openx.net.",
	"opmnstr.com.",
	"opoxv.com.",
	"oppoteammate.com.",
	"opsource.net.",
	"optad360.net.",
	"optidigital.com.",
	"optimizely.com.",
	"optimizesrv.com.",
	"optmnstr.com.",
	"optnx.com.",
	"orangeclickmedia.com.",
	"orbsrv.com.",
	"orgotitedu.info.",
	"osano.com.",
	"outbrain.com.",
	"outbrainimg.com.",
	"outlook.com.",
	"outlookmobile.com.",
	"ovvjhejceotw.com.",
	"ovvobrzryqhe.com.",
	"owneriq.net.",
	"owox.com.",
	"owrkwilxbw.com.",
	"owrxpadziajg.com.",
	"padfungusunless.com.",
	"paikoasa.tv.",
	"papmeatidigbo.com.",
	"pardot.com.",
	"parse.ly.",
	"parsely.com.",
	"partypartners.it.",
	"passport.net.",
	"payclick.it.",
	"peacto.com.",
	"pedrogarkilom.xyz.",
	"pemsrv.com.",
	"pepita.io.",
	"perfectflowing.com.",
	"perfectmarket.com.",
	"perfomail.it.",
	"perimeterx.net.",
	"permutive.app.",
	"permutive.com.",
	"petametrics.com.",
	"pghub.io.",
	"phonefactor.net.",
	"piano.io.",
	"pingdom.net.",
	"pinimg.com.",
	"pinterest.com.",
	"pippio.com.",
	"piwik.pro.",
	"pixel.wp.com.",
	"pjqchcfwtw.com.",
	"planethowl.com.",
	"platform.linkedin.com.",
	"platform.twitter.com.",
	"plausible.citynews.ovh.",
	"plausible.io.",
	"play.google.com.",
	"playfabapi.com.",
	"plug.it.",
	"pncloudfl.com.",
	"po.finance.",
	"po.st.",
	"pocket.click.",
	"pogothere.xyz.",
	"pokerstars.it.",
	"polldaddy.com.",
	"positional.ai.",
	"posthog.com.",
	"postrelease.com.",
	"powerapps.com.",
	"powerplatform.com.",
	"ppqyrngjwdq.com.",
	"prebid.org.",
	"printfriendly.com.",
	"prisasd.com.",
	"privacy-center.org.",
	"privacy-mgmt.com.",
	"privacymanager.io.",
	"pro-market.net.",
	"profitwell.com.",
	"proftrafficcounter.com.",
	"programma-affiliazione.amazon.it.",
	"proper.io.",
	"protechts.net.",
	"proxy-safebrowsing.googleapis.com.",
	"prreqcroab.icu.",
	"psichoafouts.xyz.",
	"pswec.com.",
	"ptxhzp.com.",
	"pub.network.",
	"publir.com.",
	"publytics.net.",
	"pubmatic.com.",
	"pubmnet.com.",
	"pubnation.com.",
	"pubnative.net.",
	"pubtech.ai.",
	"pubtrky.com.",
	"pulseadnetwork.com.",
	"pulsedive.com.",
	"punchh.com.",
	"pushcrew.com.",
	"pusher.com.",
	"pushio.com.",
	"px-cloud.net.",
	"qezhudifwlua.com.",
	"qsmzbddlwqoogt.com.",
	"qualaroo.com.",
	"qualtrics.com.",
	"quantcast.com.",
	"quantcount.com.",
	"quantserve.com.",
	"quantumdex.io.",
	"quantummetric.com.",
	"quora.com.",
	"rabbitrifle.com.",
	"raddoppia-bitcoin.click.",
	"rating-widget.com.",
	"ravelin.click.",
	"rcsmetrics.it.",
	"realme.com.",
	"realmemobile.com.",
	"realsrv.com.",
	"redtram.com.",
	"redventures.io.",
	"reminderlaweverything.com.",
	"repentantsympathy.com.",
	"researchnow.com.",
	"resetsrv.com.",
	"reson8.com.",
	"revjet.com.",
	"revolt.chat.",
	"rezync.com.",
	"rfihub.com.",
	"rfihub.net.",
	"richaudience.com.",
	"rkdms.com.",
	"rkgwzfwjgk.com.",
	"rkwithcatuk.org.",
	"rlcdn.com.",
	"rmshqa.com.",
	"rmtag.com.",
	"rollbar.com.",
	"rqmob.com.",
	"rqtrk.eu.",
	"rtbix.xyz.",
	"rtbsuperhub.com.",
	"rtbuzz.net.",
	"rtmark.net.",
	"rubiconproject.com.",
	"rudderlabs.com.",
	"rudderstack.com.",
	"run-syndicate.com.",
	"runative-syndicate.com.",
	"rustyanger.com.",
	"rxeosevsso.com.",
	"s-onetag.com.",
	"s3.amazonaws.com.",
	"s3xified.com.",
	"sadjklq.com.",
	"sail-horizon.com.",
	"sajari.com.",
	"salesforceliveagent.com.",
	"samplicio.us.",
	"samsungads.com.",
	"sascdn.com.",
	"sb89347.com.",
	"sbanner.com.",
	"sc-cdn.net.",
	"sc-gw.com.",
	"sc-jpl.com.",
	"sc-prod.net.",
	"sc-static.net.",
	"scambiobanner.it.",
	"scambiobanner.org.",
	"scambiobanner.tv.",
	"scambiositi.com.",
	"scarabresearch.com.",
	"scaredsnakes.com.",
	"scarletwood.com.",
	"scdn.co.",
	"scene7.com.",
	"scfsdvc.com.",
	"schip.io.",
	"scorecardresearch.com.",
	"script.ac.",
	"sda.fyi.",
	"seadform.net.",
	"securedvisit.com.",
	"securedwebark.com.",
	"securegfm.com.",
	"seedtag.com.",
	"segment.com.",
	"segment.io.",
	"semasio.net.",
	"sembox.it.",
	"sendbird.com.",
	"sentry-cdn.com.",
	"sentry.io.",
	"seolabadv.it.",
	"servedby-buysellads.com.",
	"servenobid.com.",
	"serving-sys.com.",
	"setupad.net.",
	"setupcmp.com.",
	"sexcash.com.",
	"sfbassets.com.",
	"sfx.ms.",
	"shareaholic.com.",
	"sharepointonline.com.",
	"sharethis.com.",
	"sharethrough.com.",
	"sharonwhiledemocratic.com.",
	"shiverscissors.com.",
	"sidesukbeing.org.",
	"sift.com.",
	"siftscience.com.",
	"sigheemibod.xyz.",
	"siltagefutiley.top.",
	"simpli.fi.",
	"singular.net.",
	"sitescout.com.",
	"skimresources.com.",
	"skype.com.",
	"skypeassets.com.",
	"skypeforbusiness.com.",
	"slhlfiqwekk.xyz.",
	"slicesuction.com.",
	"smaato.com.",
	"smaato.net.",
	"smartadserver.com.",
	"smartclip.net.",
	"smartclip.tv.",
	"smilewanted.com.",
	"snap.com.",
	"snapads.com.",
	"snapchat.com.",
	"snapkit.com.",
	"sneerattendingconverted.com.",
	"snoweeanalytics.com.",
	"sobisy.com.",
	"socdm.com.",
	"sojern.com.",
	"solads.media.",
	"solutionshindsight.net.",
	"sonobi.com.",
	"sophi.io.",
	"sovrn.com.",
	"spacetraff.com.",
	"speedcurve.com.",
	"spmailtechnol.com.",
	"spot.im.",
	"spotify.com.",
	"spotifycdn.com.",
	"spotim.martket.",
	"spotxchange.com.",
	"springserve.com.",
	"srvs.site.",
	"ssacdn.com.",
	"ssl-google-analytics.l.google.com.",
	"stackadapt.com.",
	"staffboozerenamed.com.",
	"staffhub.ms.",
	"starsaffiliateclub.com.",
	"startmagazine.com.",
	"startscreen.com.",
	"stat-track.com.",
	"stats.paypal.com.",
	"stats.wp.com.",
	"statsig.com.",
	"statsigapi.net.",
	"statsigapi.net.",
	"statsigcdn.com.",
	"steelhousemedia.com.",
	"stickyadstv.com.",
	"stoorgel.com.",
	"storygize.net.",
	"strigh.com.",
	"stripcdn.com.",
	"strpst.com.",
	"sumo.com.",
	"sundaysky.com.",
	"supersonicads.com.",
	"supplyframe.com.",
	"survata.com.",
	"surveymonkey.com.",
	"sway-cdn.com.",
	"sway-extensions.com.",
	"sway.com.",
	"swonqjzbc.com.",
	"syndication.twitter.com.",
	"sysnetgs.com.",
	"szbnnqyqn.com.",
	"t7cp4fldl.com.",
	"tabidmelene.com.",
	"taboola.com.",
	"taboolanews.com.",
	"tagcommander.com.",
	"tagdeliver.com.",
	"talkscreativity.com.",
	"tapad.com.",
	"tapfiliate.com.",
	"target.com.",
	"tcaabbzrack.com.",
	"tdmrfw.com.",
	"teads.tv.",
	"tealium.com.",
	"tealiumiq.com.",
	"teamblue.services.",
	"technoratimedia.com.",
	"tegleebs.com.",
	"teksishe.net.",
	"telemetry.algolia.com.",
	"temu.com.",
	"temu.to.",
	"tgtrak.com.",
	"the-ozone-project.com.",
	"thebrighttag.com.",
	"thecoreadv.com.",
	"thestartmagazine.com.",
	"threads.net.",
	"thrtle.com.",
	"tidaltv.com.",
	"tifacciounregalo.com.",
	"tinypass.com.",
	"tiqcdn.com.",
	"tncid.app.",
	"toopsoug.net.",
	"trackcmp.net.",
	"trackjs.com.",
	"trackonomics.net.",
	"trackwilltrk.com.",
	"trafficfactory.biz.",
	"trafficjunky.com.",
	"trafficjunky.net.",
	"trafficshop.com.",
	"trafostatic.com.",
	"transportgrop.live.",
	"treasuredata.com.",
	"treasurewaitlng.net.",
	"tremorhub.com.",
	"trendemon.com.",
	"tribalfusion.com.",
	"triplelift.com.",
	"troopsassistedstupidity.com.",
	"tropbikewall.art.",
	"tru.am.",
	"truoptik.com.",
	"trustarc.com.",
	"truste.com.",
	"trustedshops.com.",
	"trustradius.com.",
	"trx-hub.com.",
	"tsyndicate.com.",
	"tubemogul.com.",
	"tuobenessere.it.",
	"turboforthefirt.homes.",
	"turn.com.",
	"tvpixel.com.",
	"tvsquared.com.",
	"tynt.com.",
	"typekit.net.",
	"typography.com.",
	"tywzyhfliwdbu.com.",
	"tzegilo.com.",
	"udksgsuvcpm.com.",
	"uhqpnhorurueku.com.",
	"uncn.jp.",
	"undertone.com.",
	"unity3d.com.",
	"unrulymedia.com.",
	"unseenreport.com.",
	"uoqkbijexhfnqd.xyz.",
	"urbanairship.com.",
	"usabilla.com.",
	"usebutton.com.",
	"usefathom.com.",
	"usercentrics.com.",
	"usercentrics.eu.",
	"uservoice.com.",
	"userzoom.com.",
	"usocial.pro.",
	"utfiia.xyz.",
	"varcuringordsetts.com.",
	"vercel-analytics.com.",
	"vercel-insights.com.",
	"vhkbvpbuhwon.com.",
	"viaggiamo.it.",
	"viber.com.",
	"vidverto.io.",
	"viglink.com.",
	"vindicosuite.com.",
	"vipsohbethatlari.info.",
	"viralize.tv.",
	"virtualearth.net.",
	"visualwebsiteoptimizer.com.",
	"vivapolska.pro.",
	"vokaunget.xyz.",
	"vpixrlkggv.com.",
	"vpn-content.net.",
	"vpn-world.com.",
	"vulkanvegas.de.",
	"vungle.com.",
	"vungle.io.",
	"vuukle.com.",
	"vyfrxuytzn.com.",
	"w55c.net.",
	"waggonerfoulpillow.com.",
	"waisheph.com.",
	"wanttotrack.com.",
	"washealinginduced.com.",
	"wbbjrczne.com.",
	"wbtrk.net.",
	"wcsvjqcrpqqft.com.",
	"wdpqgagmulazv.com.",
	"weaktongue.com.",
	"webanalytics.italia.it.",
	"webclicks24.com.",
	"webcontentassessor.com.",
	"webengage.com.",
	"webflow.com.",
	"weborama.fr.",
	"webtrekk.net.",
	"widgets.wp.com.",
	"windows-ppe.net.",
	"windows.com.",
	"windows.net.",
	"windowsupdate.com.",
	"wivyiz.com.",
	"wknd.ai.",
	"worstideatum.com.",
	"wsod.com.",
	"wt-safetag.com.",
	"www-googletagmanager.l.google.com.",
	"xhylpybcad.com.",
	"ximxim.com.",
	"xlivrdr.com.",
	"xlviirdr.com.",
	"xplosion.de.",
	"xxxviijmp.com.",
	"xxxvjmp.com.",
	"yahoo.com.",
	"yahooinc.com.",
	"yahoosandbox.com.",
	"yammer.com.",
	"yammerusercontent.com.",
	"ybb-network.com.",
	"ybbserver.com.",
	"yieldlab.net.",
	"yieldlove.com.",
	"yieldmo.com.",
	"yieldoptimizer.com.",
	"ym-tack.b-cdn.net.",
	"yobee.it.",
	"yotpo.com.",
	"yottaa.com.",
	"youradexchange.com.",
	"yourwebbars.com.",
	"ysmbttmncrajnk.com.",
	"ytoworkwi.org.",
	"yvmads.com.",
	"ywhowascryin.com.",
	"zatnoh.com.",
	"zazgihgzejr.com.",
	"zemanta.com.",
	"zeusadx.com.",
	"zeustechnology.com.",
	"zi-scripts.com.",
	"zimpolo.com.",
	"zoom.us.",
	"zoominfo.com.",
	"zorosrv.com.",
	"zqtk.net.",
}
