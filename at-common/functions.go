package at_common

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	DES_KEY_HOSTPORT = "3j49d723"
)

//fileGetContents
func FileGetContents(file string) (content string, err error) {
	defer func(err *error) {
		e := recover()
		if e != nil {
			*err = fmt.Errorf("%s", e)
		}
	}(&err)
	bytes, err := ioutil.ReadFile(file)
	content = string(bytes)
	return
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
func PathExists(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
func IoBind(dst io.ReadWriter, src io.ReadWriter, fn func(err error), cfn func(count int, isPositive bool), bytesPreSec float64) {
	go func() {
		errchn := make(chan error, 2)
		go func() {
			var err error
			if bytesPreSec > 0 {
				newreader := NewReader(src)
				newreader.SetRateLimit(bytesPreSec)
				_, err = ioCopy(dst, newreader, func(c int) {
					cfn(c, false)
				})

			} else {
				_, err = ioCopy(dst, src, func(c int) {
					cfn(c, false)
				})
			}

			errchn <- err
		}()
		go func() {
			var err error
			if bytesPreSec > 0 {
				newReader := NewReader(dst)
				newReader.SetRateLimit(bytesPreSec)
				_, err = ioCopy(src, newReader, func(c int) {
					cfn(c, true)
				})
			} else {
				_, err = ioCopy(src, dst, func(c int) {
					cfn(c, true)
				})
			}
			errchn <- err
		}()
		fn(<-errchn)
	}()
}
func ioCopy(dst io.Writer, src io.Reader, fn ...func(count int)) (written int64, err error) {
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
				if len(fn) == 1 {
					fn[0](nw)
				}
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}
func TlsConnectHost(host string, timeout int) (conn tls.Conn, err error) {
	h := strings.Split(host, ":")
	port, _ := strconv.Atoi(h[1])
	return TlsConnect(h[0], port, timeout)
}

func TlsConnect(host string, port, timeout int) (conn tls.Conn, err error) {
	conf, err := getRequestTlsConfig(true)
	if err != nil {
		return
	}
	_conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Duration(timeout)*time.Millisecond)
	if err != nil {
		return
	}
	return *tls.Client(_conn, conf), err
}
func getRequestTlsConfig(isInternal bool) (conf *tls.Config, err error) {
	if isInternal {
		var cert tls.Certificate
		cert, err = tls.X509KeyPair(GetClientCert(), GetClientKey())
		if err != nil {
			return
		}
		serverCertPool := x509.NewCertPool()
		ok := serverCertPool.AppendCertsFromPEM(GetRootCert())
		if !ok {
			err = errors.New("failed to parse root certificate")
		}
		conf = &tls.Config{
			RootCAs:            serverCertPool,
			Certificates:       []tls.Certificate{cert},
			ServerName:         "anytunnel-server",
			InsecureSkipVerify: false,
		}
	} else {
		conf = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	return
}
func Connect(host string, port, timeout int) (conn net.Conn, err error) {
	conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Duration(timeout)*time.Millisecond)
	return
}
func ConnectHost(hostAndPort string, timeout int) (conn net.Conn, err error) {
	conn, err = net.DialTimeout("tcp", hostAndPort, time.Duration(timeout)*time.Millisecond)
	return
}
func ListenTls(ip string, port int) (ln *net.Listener, err error) {
	var cert tls.Certificate
	cert, err = tls.X509KeyPair(GetServerCert(), GetServerKey())
	if err != nil {
		return
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(GetRootCert())
	if !ok {
		err = errors.New("failed to parse root certificate")
	}
	config := &tls.Config{
		ClientCAs:    clientCertPool,
		ServerName:   "anytunnel-client",
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	_ln, err := tls.Listen("tcp", fmt.Sprintf("%s:%d", ip, port), config)
	if err == nil {
		ln = &_ln
	}
	return
}
func Response(channel *MessageChannel, success bool, msg string) (err error) {
	status := STATUS_FAIL
	if success {
		status = STATUS_SUCCESS
	}
	resp := MsgResponse{
		Msg:     Msg{MsgType: MSG_RESPONSE},
		Status:  status,
		Message: msg,
	}
	err = channel.Write(resp)
	return
}
func GetClusterHost(url, token, typ string) (addr string, err error) {
	if url == "" {
		err = fmt.Errorf("url is empty")
		return
	}
	if strings.Contains(url, "?") {
		url += fmt.Sprintf("&token=%s&type=%s", token, typ)
	} else {
		url += fmt.Sprintf("?token=%s&type=%s", token, typ)
	}
	d, code, err := HttpGet(url)
	if err != nil {
		return
	}
	if code != 200 {
		err = fmt.Errorf(string(d))
		return
	}
	addr = string(d)
	return
}
func GetAllInterfaceAddr() ([]net.IP, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	addresses := []net.IP{}
	for _, iface := range ifaces {

		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		// if iface.Flags&net.FlagLoopback != 0 {
		// 	continue // loopback interface
		// }
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// if ip == nil || ip.IsLoopback() {
			// 	continue
			// }
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			addresses = append(addresses, ip)
		}
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no address Found, net.InterfaceAddrs: %v", addresses)
	}
	//only need first
	return addresses, nil
}

var allChina = strings.Split("北京,海淀,东城,西城,宣武,丰台,朝阳,崇文,大兴,石景山,门头沟,房山,通州,顺义,怀柔,昌平,平谷,密云县,延庆县,天津,和平,河西,河北,河东,南开,红桥,北辰,津南,武清,塘沽,西青,汉沽,大港,宝坻,东丽,蓟县,静海县,宁河县,上海,黄浦,卢湾,徐汇,长宁,静安,普陀,闸北,杨浦,虹口,闵行,宝山,嘉定,浦东新,金山,松江,青浦,南汇,奉贤,崇明县,重庆,渝中,大渡口,江北,沙坪坝,九龙坡,南岸,北碚,万盛,双桥,渝北,巴南,万州,涪陵,黔江,长寿,江津,永川,南川,綦江县,潼南县,铜梁县,大足县,荣昌县,璧山县,垫江县,武隆县,丰都县,城口县,梁平县,开县,巫溪县,巫山县,奉节县,云阳县,忠县,石柱土家族,彭水苗族土家族,酉阳苗族,秀山土家族苗族,新疆维吾尔,乌鲁木齐,克拉玛依,石河子,阿拉尔,图木舒克,五家渠,哈密,吐鲁番,阿克苏,喀什,和田,伊宁,塔城,阿勒泰,奎屯,博乐,昌吉,阜康,库尔勒,阿图什,乌苏,西藏,拉萨,日喀则,宁夏回族,银川,石嘴山,吴忠,固原,中卫,青,铜峡,灵武,内蒙古,呼和浩特,包头,乌海,赤峰,通辽,鄂尔多斯,呼伦贝尔,巴彦淖尔,乌兰察布,霍林郭勒,满洲里,牙克石,扎兰屯,根河,额尔古纳,丰镇,锡林浩特,二连浩特,乌兰浩特,阿尔山,广西壮族,南宁,柳州,桂林,梧州,北海,崇左,来宾,贺州,玉林,百色,河池,钦州,防城港,贵港,岑溪,凭祥,合山,北流,宜州,东兴,桂平,黑龙江,哈尔滨,大庆,齐齐哈尔,佳木斯,鸡西,鹤岗,双鸭山,牡丹江,伊春,七台河,黑河,绥化,五常,双城,尚志,纳河,虎林,密山,铁力,同江,富锦,绥芬河,海林,宁安,穆林,北安,五大连池,肇东,海伦,安达,长春,吉林,四平,辽源,通化,白山,松原,白城,九台,榆树,德惠,舒兰,桦甸,蛟河,磐石,公主岭,双辽,梅河口,集安,临江,大安,洮南,延吉,图们,敦化,龙井,珲春,和龙,辽宁,沈阳,大连,鞍山,抚顺,本溪,丹东,锦州,营口,阜新,辽阳,盘锦,铁岭,朝阳,葫芦岛,新民,瓦房店,普兰,庄河,海城,东港,凤城,凌海,北镇,大石桥,盖州,灯塔,调兵山,开原,凌源,北票,兴城,河北,石家庄,唐山,邯郸,秦皇岛,保定,张家口,承德,廊坊,沧州,衡水,邢台,辛集,藁城,晋州,新乐,鹿泉,遵化,迁安,武安,南宫,沙河,涿州,定州,安国,高碑店,泊头,任丘,黄骅,河间,霸州,三河,冀州,深州,山东,济南,青岛,淄博,枣庄,东营,烟台,潍坊,济宁,泰安,威海,日照,莱芜,临沂,德州,聊城,菏泽,滨州,章丘,胶南,胶州,平度,莱西,即墨,滕州,龙口,莱阳,莱州,招远,蓬莱,栖霞,海阳,青州,诸城,安丘,高密,昌邑,兖州,曲阜,邹城,乳山,文登,荣成,乐陵,临清,禹城,江苏,南京,镇江,常州,无锡,苏州,徐州,连云港,淮安,盐城,扬州,泰州,南通,宿迁,江阴,宜兴,邳州,新沂,金坛,溧阳,常熟,张家港,太仓,昆山,吴江,如皋,通州,海门,启东,东台,大丰,高邮,江都,仪征,丹阳,扬中,句容,泰兴,姜堰,靖江,兴化,安徽,合肥,蚌埠,芜湖,淮南,亳州,阜阳,淮北,宿州,滁州,安庆,巢湖,马鞍山,宣城,黄山,池州,铜陵,界首,天长,明光,桐城,宁国,浙江,杭州,嘉兴,湖州,宁波,金华,温州,丽水,绍兴,衢州,舟山,台州,建德,富阳,临安,余姚,慈溪,奉化,瑞安,乐清,海宁,平湖,桐乡,诸暨,上虞,嵊州,兰溪,义乌,东阳,永康,江山,临海,温岭,龙泉,福建,福州,厦门,泉州,三明,南平,漳州,莆田,宁德,龙岩,福清,长乐,永安,石狮,晋江,南安,龙海,邵武,武夷山,建瓯,建阳,漳平,福安,福鼎,广东,广州,深圳,汕头,惠州,珠海,揭阳,佛山,河源,阳江,茂名,湛江,梅州,肇庆,韶关,潮州,东莞,中山,清远,江门,汕尾,云浮,增城,从化,乐昌,南雄,台山,开平,鹤山,恩平,廉江,雷州,吴川,高州,化州,高要,四会,兴宁,陆丰,阳春,英德,连州,普宁,罗定,海南,海口,三亚,琼海,文昌,万宁,五指山,儋州,东方,云南,昆明,曲靖,玉溪,保山,昭通,丽江,普洱,临沧,安宁,宣威,个旧,开远,景洪,楚雄,大理,潞西,瑞丽,贵州,贵阳,六盘水,遵义,安顺,清镇,赤水,仁怀,铜仁,毕节,兴义,凯里,都匀,福泉,四川,成都,绵阳,德阳,广元,自贡,攀枝花,乐山,南充,内江,遂宁,广安,泸州,达州,眉山,宜宾,雅安,资阳,都江堰,彭州,邛崃,崇州,广汉,什邡,绵竹,江油,峨眉山,阆中,华蓥,万源,简阳,西昌,湖南,长沙,株洲,湘潭,衡阳,岳阳,郴州,永州,邵阳,怀化,常德,益阳,张家界,娄底,浏阳,醴陵,湘乡,韶山,耒阳,常宁,武冈,临湘,汨罗津,沅江,资兴,洪江,冷水江,涟源,吉首,湖北,武汉,襄樊,宜昌,黄石,鄂州,随州,荆州,荆门,十堰,孝感,黄冈,咸宁,大冶,丹江口,洪湖,石首,松滋,宜都,当阳,枝江,老河口,枣阳,宜城,钟祥,应城,安陆,汉川,麻城,武穴,赤壁,广水,仙桃,天门,潜江,恩施,利川,河南,郑州,洛阳,开封,漯河,安阳,新乡,周口,三门峡,焦作,平顶山,信阳,南阳,鹤壁,濮阳,许昌,商丘,驻马店,巩义,新郑,新密,登封,荥阳,偃师,汝州,舞钢,林州,卫辉,辉县,沁阳,孟州,禹州,长葛,义马,灵宝,邓州,永城,项城,济源,山西,太原,大同,忻州,阳泉,长治,晋城,朔州,晋中,运城,临汾,吕梁,古交,潞城,高平,介休,永济,河津,原平,侯马,霍州,孝义,汾阳,陕西,西安,咸阳,铜川,延安,宝鸡,渭南,汉中,安康,商洛,榆林,兴平,韩城,华阴,甘肃,兰州,天水,平凉,酒泉,嘉峪关,金昌,白银,武威,张掖,庆阳,定西,陇南,玉门,敦煌,临夏,合作,青海,西宁,格尔木,德令哈,江西,南昌,九江,赣州,吉安,鹰潭,上饶,萍乡,景德镇,新余,宜春,抚州,乐平,瑞昌,贵溪,瑞金,南康,井冈山,丰城,樟树,高安,德兴", ",")

func IsChina(country string) bool {
	for _, v := range allChina {
		if strings.HasPrefix(country, v) {
			return true
		}
	}
	return false
}

//md5加密
func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}
func UDPPacket(srcAddr string, packet []byte) []byte {
	addrBytes := []byte(srcAddr)
	addrLength := uint16(len(addrBytes))
	bodyLength := uint16(len(packet))
	pkg := new(bytes.Buffer)
	binary.Write(pkg, binary.LittleEndian, addrLength)
	binary.Write(pkg, binary.LittleEndian, addrBytes)
	binary.Write(pkg, binary.LittleEndian, bodyLength)
	binary.Write(pkg, binary.LittleEndian, packet)
	return pkg.Bytes()
}
func ReadUDPPacket(conn *tls.Conn) (srcAddr string, packet []byte, err error) {
	reader := bufio.NewReader(conn)
	var addrLength uint16
	var bodyLength uint16
	err = binary.Read(reader, binary.LittleEndian, &addrLength)
	if err != nil {
		return
	}
	_srcAddr := make([]byte, addrLength)
	n, err := reader.Read(_srcAddr)
	if err != nil {
		return
	}
	if n != int(addrLength) {
		return
	}
	srcAddr = string(_srcAddr)

	err = binary.Read(reader, binary.LittleEndian, &bodyLength)
	if err != nil {
		return
	}
	packet = make([]byte, bodyLength)
	n, err = reader.Read(packet)
	if err != nil {
		return
	}
	if n != int(bodyLength) {
		return
	}
	return
}
