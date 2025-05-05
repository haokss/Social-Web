package report

import (
	"fmt"
	"strconv"
	"strings"
	"todo_list/model"

	"github.com/xuri/excelize/v2"
)

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

		amount, _ := strconv.ParseFloat(data["金额（元）"], 64)
		result = append(result, model.Bill{
			TransactionTime: data["交易时间"],
			TransactionType: data["交易类型"],
			Counterparty:    data["交易对方"],
			Product:         data["商品"],
			IncomeExpense:   data["收/支"],
			Amount:          amount,
			PaymentMethod:   data["支付方式"],
			Status:          data["当前状态"],
			TransactionID:   data["交易单号"],
			MerchantID:      data["商户单号"],
			Remarks:         data["备注"],
		})
	}

	return result, nil
}

func parseAlipayBill(path string) ([]map[string]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}

	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	var result []map[string]string
	startIndex := 0
	for i, row := range rows {
		if len(row) > 0 && strings.Contains(row[0], "交易号") {
			startIndex = i
			break
		}
	}

	if startIndex == 0 {
		return nil, fmt.Errorf("未找到表头")
	}

	headers := rows[startIndex]
	for _, row := range rows[startIndex+1:] {
		record := map[string]string{}
		for i, col := range row {
			if i < len(headers) {
				record[headers[i]] = strings.TrimSpace(col)
			}
		}
		result = append(result, record)
	}

	return result, nil
}
