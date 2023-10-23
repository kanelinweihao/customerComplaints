package dataGet

import (
	"fmt"
	"go.lwh.com/linweihao/customerComplaints/factory/factoryOfDB"
	"go.lwh.com/linweihao/customerComplaints/factory/factoryOfGoroutine"
	"go.lwh.com/linweihao/customerComplaints/utils/db"
	"go.lwh.com/linweihao/customerComplaints/utils/err"
	"go.lwh.com/linweihao/customerComplaints/utils/goroutine"
	"go.lwh.com/linweihao/customerComplaints/utils/rfl"
	"go.lwh.com/linweihao/customerComplaints/utils/time"
)

type TypeEntityData interface {
	Entity1 | Entity2 | Entity3 | Entity4
}

var query1 string = `SELECT
CONCAT('北京时间 ',DATE_SUB(O.created_at, INTERVAL -8 HOUR)) AS '购买时间(北京时间)',
O.user_id_buy AS 'UID',
CONCAT('手机号 ',U.mobile_no) AS '手机号',
P.sku_id AS '藏品集编号',
P.product_name_cn AS '商品名称',
P.product_fiat_original AS '价格单位',
P.product_price_original AS '商品价格',
CONCAT('订单号 ',O.product_order_no) AS '订单号',
N.goods_no AS '藏品唯一编号',
N.user_id AS '藏品当前持有人UID',
IF(O.user_id_buy = N.user_id,'仍然持有','未持有') AS '是否持有'
FROM
(
SELECT
id,
created_at,
user_id_buy,
product_order_no,
product_id,
REPLACE(JSON_EXTRACT(product_order_config, '$.goods_no'), '"', '') AS 'goods_no'
FROM product_order_pool_nft
WHERE status = 5
AND org_id = 33
AND user_id_buy = ?
) AS O
LEFT JOIN product_pool_nft AS P
ON O.product_id  = P.id
LEFT JOIN users AS U
ON O.user_id_buy = U.id
INNER JOIN nft_goods AS N
ON O.goods_no = N.goods_no
ORDER BY O.user_id_buy DESC, O.id DESC
;
`

var query2 string = `SELECT
CONCAT('北京时间 ',DATE_SUB(O.created_at, INTERVAL -8 HOUR)) AS '成交时间(北京时间)',
O.user_id_sell AS '卖方UID',
CONCAT('手机号 ',US.mobile_no) AS '卖方手机号',
O.user_id_buy AS '买方UID',
CONCAT('手机号 ',UB.mobile_no) AS '买方手机号',
P.sku_id AS '藏品集编号',
P.product_name_cn AS '藏品名称',
C.code_id AS '铸造编号',
P.price_original AS '成交价格',
CONCAT('订单单号 ',O.product_order_no) AS '订单单号'
FROM product_order_nft AS O
LEFT JOIN product_nft AS P
ON O.product_id  = P.id
LEFT JOIN nft_goods AS N
ON P.product_name_salt = N.goods_no
LEFT JOIN cast_item_token AS C
ON N.sku_id = C.sku_id
AND N.seq_index = C.seq_index
LEFT JOIN users AS UB
ON O.user_id_buy = UB.id
LEFT JOIN users AS US
ON O.user_id_sell = US.id
WHERE O.status = 5
AND O.org_id = 33
AND O.user_id_buy = ?
ORDER BY O.user_id_buy DESC, O.id DESC
;
`

var query3 string = `SELECT
CONCAT('北京时间 ',DATE_SUB(O.created_at, INTERVAL -8 HOUR)) AS '成交时间(北京时间)',
O.user_id_sell AS '卖方UID',
CONCAT('手机号 ',US.mobile_no) AS '卖方手机号',
O.user_id_buy AS '买方UID',
CONCAT('手机号 ',UB.mobile_no) AS '买方手机号',
P.sku_id AS '藏品集编号',
P.product_name_cn AS '藏品名称',
C.code_id AS '铸造编号',
P.price_original AS '成交价格',
CONCAT('订单单号 ',O.product_order_no) AS '订单单号'
FROM product_order_nft AS O
LEFT JOIN product_nft AS P
ON O.product_id  = P.id
LEFT JOIN nft_goods AS N
ON P.product_name_salt = N.goods_no
LEFT JOIN cast_item_token AS C
ON N.sku_id = C.sku_id
AND N.seq_index = C.seq_index
LEFT JOIN users AS UB
ON O.user_id_buy = UB.id
LEFT JOIN users AS US
ON O.user_id_sell = US.id
WHERE O.status = 5
AND O.org_id = 33
AND O.user_id_sell = ?
ORDER BY O.user_id_sell DESC, O.id DESC
;
`

var query4 string = `SELECT
A.user_id AS 'UID',
A.sum_amount AS '可领取水晶总数',
IFNULL(B.sum_amount,0) AS '已领取水晶总数'
FROM (
SELECT
user_id,
SUM(amount) AS 'sum_amount'
FROM user_crystal_log
GROUP BY user_id
) AS A
LEFT JOIN (
SELECT
user_id,
SUM(amount) AS 'sum_amount'
FROM user_crystal_log
WHERE type = 1
GROUP BY user_id
) AS B
ON A.user_id = B.user_id
WHERE A.user_id = ?
ORDER BY A.user_id DESC
;
`

type Entity1 struct {
	BuyAt          string `db:"购买时间(北京时间)"`
	UserId         int    `db:"UID"`
	MobileNo       string `db:"手机号"`
	SkuId          int    `db:"藏品集编号"`
	Title          string `db:"商品名称"`
	PriceUnit      string `db:"价格单位"`
	Price          string `db:"商品价格"`
	ProductOrderNo string `db:"订单号"`
	GoodsNo        string `db:"藏品唯一编号"`
	UserIdHold     int    `db:"藏品当前持有人UID"`
	MsgHold        string `db:"是否持有"`
}

type Entity2 struct {
	BuyAt          string `db:"成交时间(北京时间)"`
	UserIdSell     string `db:"卖方UID"`
	MobileNoSell   string `db:"卖方手机号"`
	UserIdBuy      string `db:"买方UID"`
	MobileNoBuy    string `db:"买方手机号"`
	SkuId          string `db:"藏品集编号"`
	Title          string `db:"藏品名称"`
	CodeId         string `db:"铸造编号"`
	Price          string `db:"成交价格"`
	ProductOrderNo string `db:"订单单号"`
}

type Entity3 struct {
	BuyAt          string `db:"成交时间(北京时间)"`
	UserIdSell     string `db:"卖方UID"`
	MobileNoSell   string `db:"卖方手机号"`
	UserIdBuy      string `db:"买方UID"`
	MobileNoBuy    string `db:"买方手机号"`
	SkuId          string `db:"藏品集编号"`
	Title          string `db:"藏品名称"`
	CodeId         string `db:"铸造编号"`
	Price          string `db:"成交价格"`
	ProductOrderNo string `db:"订单单号"`
}

type Entity4 struct {
	UserId          string `db:"UID"`
	SumAmountCanGet string `db:"可领取水晶总数"`
	SumAmountGeted  string `db:"已领取水晶总数"`
}

var Params1 = rfl.ArrParams{
	"BuyAt": {
		"field": "BuyAt",
		"title": "购买时间(北京时间)",
		"sort":  1,
	},
	"UserId": {
		"field": "UserId",
		"title": "UID",
		"sort":  2,
	},
	"MobileNo": {
		"field": "MobileNo",
		"title": "手机号",
		"sort":  3,
	},
	"SkuId": {
		"field": "SkuId",
		"title": "藏品集编号",
		"sort":  4,
	},
	"Title": {
		"field": "Title",
		"title": "商品名称",
		"sort":  5,
	},
	"PriceUnit": {
		"field": "PriceUnit",
		"title": "价格单位",
		"sort":  6,
	},
	"Price": {
		"field": "Price",
		"title": "商品价格",
		"sort":  7,
	},
	"ProductOrderNo": {
		"field": "ProductOrderNo",
		"title": "订单号",
		"sort":  8,
	},
	"GoodsNo": {
		"field": "GoodsNo",
		"title": "藏品唯一编号",
		"sort":  9,
	},
	"UserIdHold": {
		"field": "UserIdHold",
		"title": "藏品当前持有人UID",
		"sort":  10,
	},
	"MsgHold": {
		"field": "MsgHold",
		"title": "是否持有",
		"sort":  11,
	},
}

var Params2 = rfl.ArrParams{
	"BuyAt": {
		"field": "BuyAt",
		"title": "成交时间(北京时间)",
		"sort":  1,
	},
	"UserIdSell": {
		"field": "UserIdSell",
		"title": "卖方UID",
		"sort":  2,
	},
	"MobileNoSell": {
		"field": "MobileNoSell",
		"title": "卖方手机号",
		"sort":  3,
	},
	"UserIdBuy": {
		"field": "UserIdBuy",
		"title": "买方UID",
		"sort":  4,
	},
	"MobileNoBuy": {
		"field": "MobileNoBuy",
		"title": "买方手机号",
		"sort":  5,
	},
	"SkuId": {
		"field": "SkuId",
		"title": "藏品集编号",
		"sort":  6,
	},
	"Title": {
		"field": "Title",
		"title": "藏品名称",
		"sort":  7,
	},
	"CodeId": {
		"field": "CodeId",
		"title": "铸造编号",
		"sort":  8,
	},
	"Price": {
		"field": "Price",
		"title": "成交价格",
		"sort":  9,
	},
	"ProductOrderNo": {
		"field": "ProductOrderNo",
		"title": "订单单号",
		"sort":  10,
	},
}

var Params3 = rfl.ArrParams{
	"BuyAt": {
		"field": "BuyAt",
		"title": "成交时间(北京时间)",
		"sort":  1,
	},
	"UserIdSell": {
		"field": "UserIdSell",
		"title": "卖方UID",
		"sort":  2,
	},
	"MobileNoSell": {
		"field": "MobileNoSell",
		"title": "卖方手机号",
		"sort":  3,
	},
	"UserIdBuy": {
		"field": "UserIdBuy",
		"title": "买方UID",
		"sort":  4,
	},
	"MobileNoBuy": {
		"field": "MobileNoBuy",
		"title": "买方手机号",
		"sort":  5,
	},
	"SkuId": {
		"field": "SkuId",
		"title": "藏品集编号",
		"sort":  6,
	},
	"Title": {
		"field": "Title",
		"title": "藏品名称",
		"sort":  7,
	},
	"CodeId": {
		"field": "CodeId",
		"title": "铸造编号",
		"sort":  8,
	},
	"Price": {
		"field": "Price",
		"title": "成交价格",
		"sort":  9,
	},
	"ProductOrderNo": {
		"field": "ProductOrderNo",
		"title": "订单单号",
		"sort":  10,
	},
}

var Params4 = rfl.ArrParams{
	"UserId": {
		"field": "UserId",
		"title": "UID",
		"sort":  1,
	},
	"SumAmountCanGet": {
		"field": "SumAmountCanGet",
		"title": "可领取水晶总数",
		"sort":  2,
	},
	"SumAmountGeted": {
		"field": "SumAmountGeted",
		"title": "已领取水晶总数",
		"sort":  3,
	},
}

/*
func
*/

func GetDataFromDB(userId int) (boxData rfl.BoxData, boxParams rfl.BoxParams[int]) {
	// time.ShowTimeAndMsg("Data get begin")
	// entityDB := db.InitDB()
	entityDB := factoryOfDB.MakeEntityOfDB()
	defer entityDB.CloseDB()
	// init
	boxData = rfl.BoxData{}
	boxParams = rfl.BoxParams[int]{}
	// params
	boxParams[1] = Params1
	boxParams[2] = Params2
	boxParams[3] = Params3
	boxParams[4] = Params4
	// write
	// 1
	var arrEntity1 []Entity1
	channelRead1 := getArrAttrForExcelOfWrite(
		entityDB,
		arrEntity1,
		query1,
		userId,
		1)
	// 2
	var arrEntity2 []Entity2
	channelRead2 := getArrAttrForExcelOfWrite(
		entityDB,
		arrEntity2,
		query2,
		userId,
		2)
	// 3
	var arrEntity3 []Entity3
	channelRead3 := getArrAttrForExcelOfWrite(
		entityDB,
		arrEntity3,
		query3,
		userId,
		3)
	// 4
	var arrEntity4 []Entity4
	channelRead4 := getArrAttrForExcelOfWrite(
		entityDB,
		arrEntity4,
		query4,
		userId,
		4)
	// read
	arrAttrForExcel1 := getArrAttrForExcelOfRead(
		channelRead1,
		1)
	arrAttrForExcel2 := getArrAttrForExcelOfRead(
		channelRead2,
		2)
	arrAttrForExcel3 := getArrAttrForExcelOfRead(
		channelRead3,
		3)
	arrAttrForExcel4 := getArrAttrForExcelOfRead(
		channelRead4,
		4)
	boxData[1] = arrAttrForExcel1
	boxData[2] = arrAttrForExcel2
	boxData[3] = arrAttrForExcel3
	boxData[4] = arrAttrForExcel4
	// return
	time.ShowTimeAndMsg("Data get success")
	return boxData, boxParams
}

// func getArrAttrForExcelOfWrite[T Entity1|Entity2|Entity3|Entity4] (entityDB *db.EntityDB, arrEntityNum []T, queryNum string, userId int, i int) (entityChannelNum *goroutine.EntityChannel) {
func getArrAttrForExcelOfWrite[T TypeEntityData](entityDB *db.EntityDB, arrEntityNum []T, queryNum string, userId int, i int) (entityChannelNum *goroutine.EntityChannel) {
	// entityChannelNum = goroutine.InitEntityChannel()
	entityChannelNum = factoryOfGoroutine.MakeEntityOfGoroutine()
	go db.GetArrAttrForExcelUseGoroutine(
		entityDB,
		entityChannelNum,
		arrEntityNum,
		queryNum,
		userId)
	return entityChannelNum
}

func getArrAttrForExcelOfRead(entityChannelNum *goroutine.EntityChannel, i int) (arrAttrForExcelNum rfl.ArrAttrForExcel) {
	dataOnce := entityChannelNum.ReadOnce()
	arrAttrForExcelNum, ok := dataOnce.(rfl.ArrAttrForExcel)
	if !ok {
		strNum := rfl.IntToStr(i)
		msgError := fmt.Sprintf(
			"The type of |%s| is not |%s|",
			"arrAttrForExcel"+strNum,
			"rfl.ArrAttrForExcel")
		err.ErrPanic(msgError)
	}
	return arrAttrForExcelNum
}
