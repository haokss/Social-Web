package report

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"
	"todo_list/model"

	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 解析微信xlsx账单
func ParseWeChatXLSX(path string) ([]model.Bill, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	// 表头一般在第17行，数据从第18行开始
	headers := rows[16]
	var result []model.Bill

	for _, row := range rows[17:] {
		data := map[string]string{}
		for i, cell := range row {
			if i < len(headers) {
				data[headers[i]] = strings.TrimSpace(cell)
			}
		}
		rawAmount := strings.TrimSpace(data["金额(元)"])
		rawAmount = strings.TrimPrefix(rawAmount, "￥")
		amount, _ := strconv.ParseFloat(rawAmount, 64)
		// 转化导入时间格式
		layoutInput := "1/2/06 15:04"
		layoutOutput := "2006-01-02 15:04:05"
		rawTime := strings.TrimSpace(data["交易时间"])
		var formattedTime string
		if parsedTime, err := time.Parse(layoutInput, rawTime); err == nil {
			formattedTime = parsedTime.Format(layoutOutput)
		} else {
			// fmt.Printf("解析时间失败: %v, 原始值: %s\n", err, rawTime)
			formattedTime = rawTime // 或者直接 continue 跳过该行
		}
		result = append(result, model.Bill{
			TransactionTime: formattedTime,
			TransactionType: data["交易类型"],
			Counterparty:    data["交易对方"],
			Product:         data["商品"],
			IncomeExpense:   data["收/支"],
			Amount:          amount,
			PaymentMethod:   "微信支付",
			Status:          data["当前状态"],
			TransactionID:   data["交易单号"],
			MerchantID:      data["商户单号"],
			Remarks:         data["备注"],
		})
	}

	return result, nil
}

// 解析阿里账单csv格式
func ParseAlipayCSV(path string) ([]model.Bill, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 修改编码方式解码 GBK -> UTF-8
	utf8Reader := transform.NewReader(file, simplifiedchinese.GBK.NewDecoder())

	reader := csv.NewReader(utf8Reader)
	reader.FieldsPerRecord = -1 // 允许变长行

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// 找到表头
	startIndex := 4
	headers := rows[startIndex]
	var bills []model.Bill

	// 清理 headers 空格
	for i := range headers {
		headers[i] = strings.TrimSpace(headers[i])
	}

	for _, row := range rows[startIndex+1:] {
		if len(row) == 0 || row[0] == "" {
			continue
		}
		// 遇到结束标志就停止解析
		if strings.Contains(row[0], "----------------") {
			break
		}
		data := map[string]string{}
		for i := range headers {
			if i < len(row) {
				data[headers[i]] = strings.TrimSpace(row[i])
			}
		}

		// 解析金额（可能为负数）
		amountStr := strings.ReplaceAll(data["金额（元）"], ",", "")
		amount, _ := strconv.ParseFloat(amountStr, 64)

		bill := model.Bill{
			TransactionTime: data["付款时间"],
			TransactionType: data["类型"],
			Counterparty:    data["交易对方"],
			Product:         data["商品名称"],
			IncomeExpense:   data["收/支"],
			Amount:          amount,
			PaymentMethod:   "支付宝",
			Status:          data["交易状态"],
			TransactionID:   data["交易号"],
			MerchantID:      data["商家订单号"],
			Remarks:         data["备注"],
		}
		bills = append(bills, bill)
	}

	return bills, nil
}
