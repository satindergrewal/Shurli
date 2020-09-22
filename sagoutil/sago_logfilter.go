package sagoutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SwapStatus defines the data type to store filtered data and push to application in JSON format for UI side rendering.
type SwapStatus struct {
	State        string  `json:"state,omitempty"`
	StateID      string  `json:"state_id,omitempty"`
	Status       string  `json:"status,omitempty"`
	StateHash    string  `json:"state_hash,omitempty"`
	BaseTxID     string  `json:"base_txid,omitempty"`
	RelTxID      string  `json:"rel_txid,omitempty"`
	SwapID       string  `json:"swap_id,omitempty"`
	SendToAddr   string  `json:"sendtoaddr,omitempty"`
	RecvAddr     string  `json:"recvaddr,omitempty"`
	Base         string  `json:"base,omitempty"`
	Rel          string  `json:"rel,omitempty"`
	BaseAmount   float64 `json:"base_amount,omitempty"`
	RelAmount    float64 `json:"rel_amount,omitempty"`
	SwapFullData string  `json:"swap_full_data,omitempty"`
}

//ZFrom to render the JSON data from "from." log entery coming from subatomic stdout
type ZFrom []struct {
	Address string  `json:"address,omitempty"`
	Amount  float64 `json:"amount,omitempty"`
	Memo    string  `json:"memo,omitempty"`
}

// SwapHistory type stores the single object from Swaps logs history
type SwapHistory struct {
	SwapID           string       `json:"swapid"`
	TimeStamp        string       `json:"timestamp"`
	Base             string       `json:"base"`
	Rel              string       `json:"rel"`
	BaseAmount       float64      `json:"baseamount"`
	RelAmount        float64      `json:"relamount"`
	BaseTxID         string       `json:"basetxid"`
	RelTxID          string       `json:"reltxid"`
	Status           string       `json:"status"`
	BobID            string       `json:"bobid"`
	BobPubkey        string       `json:"bobpubkey"`
	BobAuthorized    bool         `json:"bobauthorized"`
	OrderIDTimeStamp string       `json:"orderidtimestamp"`
	OrderPrice       float64      `json:"orderprice"`
	OrderMaxVolume   float64      `json:"ordermaxvolume"`
	BaseIcon         string       `json:"baseicon"`
	RelIcon          string       `json:"relicon"`
	BaseExplorer     string       `json:"baseexplorer"`
	RelExplorer      string       `json:"relexplorer"`
	ZBase            bool         `json:"zbase"`
	ZRel             bool         `json:"zrel"`
	SwapLog          []SwapStatus `json:"swaplog,omitempty"`
}

// SwapsHistory type is used as collection of SwapHistory objects
type SwapsHistory []SwapHistory

// SwapsHistory returns processed slice of swaplogs in JSON format
func (history SwapsHistory) SwapsHistory() (SwapsHistory, error) {

	files, err := ioutil.ReadDir("swaplogs")
	if err != nil {
		log.Println(err)
		return history, errors.New(err.Error())
	}

	// var multipleLogs = []string{logString, logString0, logString1}
	for _, file := range files {
		// fmt.Println(i)
		// fmt.Println(file)

		fn := strings.Split(file.Name(), "_")
		fnid := strings.Split(fn[1], ".")
		timestamp := fn[0]
		swapid := fnid[0]
		// fmt.Printf("\ntimestamp: %s\nswapid: %s\n", fn, fnid)

		fileRead, err := ioutil.ReadFile("swaplogs/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}

		logval, _ := SwapLogFilter(string(fileRead), "full")
		// fmt.Println("logval", logval)
		var _logval []SwapStatus
		err = json.Unmarshal([]byte(logval), &_logval)
		if err != nil {
			log.Println(err)
		}

		// fmt.Println(len(_logval))
		// fmt.Println(_logval[6].RelTxID)

		var reltxid, basetxid, swapstatus string
		swapstatus = "INCOMPLETE"
		for _, v := range _logval {
			// fmt.Println(v.State)
			// fmt.Println(v.Status)
			if v.RelTxID != "" {
				reltxid = v.RelTxID
			}
			if v.BaseTxID != "" {
				basetxid = v.BaseTxID
			}
			if v.State == "SWAP COMPLETE" {
				swapstatus = "COMPLETE"
			}
		}

		var orderIDInfo []interface{}
		json.Unmarshal([]byte(_logval[len(_logval)-1].SwapFullData), &orderIDInfo)

		// if err != nil {
		// 	fmt.Println(err)
		// }

		// fmt.Println(orderIDInfo[5])
		// fmt.Println("authorized: ", orderIDInfo[7].(map[string]interface{})["authorized"])

		baseAmount, err := strconv.ParseFloat(fmt.Sprintf("%v", orderIDInfo[3]), 64)
		relAmount, err := strconv.ParseFloat(fmt.Sprintf("%v", orderIDInfo[4]), 64)
		authorized, err := strconv.ParseBool(fmt.Sprintf("%v", orderIDInfo[7].(map[string]interface{})["authorized"]))
		zbase, err := strconv.ParseBool(fmt.Sprintf("%v", orderIDInfo[7].(map[string]interface{})["zbase"]))
		zrel, err := strconv.ParseBool(fmt.Sprintf("%v", orderIDInfo[7].(map[string]interface{})["zrel"]))
		price, err := strconv.ParseFloat(fmt.Sprintf("%v", orderIDInfo[7].(map[string]interface{})["price"]), 64)
		maxvolume, err := strconv.ParseFloat(fmt.Sprintf("%v", orderIDInfo[7].(map[string]interface{})["max_volume"]), 64)

		timestampint, err := strconv.ParseInt(timestamp, 10, 64)
		var unixTime int64 = timestampint
		t := time.Unix(unixTime, 0)
		// fmt.Println("t", t)
		strDate := t.Format(time.UnixDate)
		// fmt.Println(strDate)

		history = append(history, SwapHistory{
			SwapID:           swapid,
			TimeStamp:        strDate,
			Base:             fmt.Sprintf("%v", orderIDInfo[7].(map[string]interface{})["base"]),
			Rel:              fmt.Sprintf("%v", orderIDInfo[7].(map[string]interface{})["rel"]),
			BaseAmount:       baseAmount,
			RelAmount:        relAmount,
			BaseTxID:         basetxid,
			RelTxID:          reltxid,
			BobID:            fmt.Sprintf("%v", orderIDInfo[7].(map[string]interface{})["handle"]),
			BobPubkey:        fmt.Sprintf("%v", orderIDInfo[7].(map[string]interface{})["pubkey"]),
			BobAuthorized:    authorized,
			OrderIDTimeStamp: fmt.Sprintf("%v", orderIDInfo[7].(map[string]interface{})["timestamp"]),
			OrderPrice:       price,
			OrderMaxVolume:   maxvolume,
			BaseIcon:         fmt.Sprintf("%v", orderIDInfo[7].(map[string]interface{})["baseicon"]),
			RelIcon:          fmt.Sprintf("%v", orderIDInfo[7].(map[string]interface{})["relicon"]),
			BaseExplorer:     fmt.Sprintf("%v", orderIDInfo[5]),
			RelExplorer:      fmt.Sprintf("%v", orderIDInfo[6]),
			ZBase:            zbase,
			ZRel:             zrel,
			Status:           swapstatus,
			// SwapLog:    _logval,
		})
		// fmt.Println("\n", history)
	}

	// var historyJSON []byte
	// // historyJSON, _ = json.Marshal(history)
	// // fmt.Println(len(history))
	// historyJSON, _ = json.MarshalIndent(history, "", "  ")
	// fmt.Println(string(historyJSON))
	// return string(historyJSON)
	return history, nil
}

// SwapLogFilter returns JSON processed data for submitted log string
func SwapLogFilter(logString, answer string) (string, error) {
	var expOpenReq = regexp.MustCompile(`(?m)openrequest\..+$`)
	openReq := expOpenReq.FindString(logString)
	// fmt.Println(openReq)
	openReqSf := strings.Fields(openReq)
	// fmt.Println(openReqSf[2])

	var statuses []SwapStatus

	if len(openReqSf) > 0 {
		// fmt.Printf("length of openReqSf is greater: %d\n", len(openReqSf))

		openReqHash := strings.ReplaceAll(openReqSf[2], "(", "")
		openReqHash = strings.ReplaceAll(openReqHash, ")", "")
		// fmt.Println(openReqHash)
		openReqID := strings.Split(openReqSf[0], ".")
		// fmt.Println(openReqID)
		// fmt.Println(openReqID[0])
		// fmt.Println(openReqID[1])
		// fmt.Println("Channel Open Request Sent:", openReqID[1])
		state0 := SwapStatus{
			Status:    "0",
			State:     openReqID[0],
			StateID:   openReqID[1],
			StateHash: openReqHash,
		}

		// fmt.Println("State 0:", state0)
		state0JSON, _ := json.Marshal(state0)
		// fmt.Println("Channel Opened")
		// fmt.Println("state0 JSON:", string(state0JSON))
		if answer == "single" {
			return string(state0JSON), nil
		}
		if answer == "full" {
			statuses = append(statuses, state0)
		}

	} //else {
	// fmt.Printf("length of openReqSf is lower: %d\n", len(openReqSf))
	//}

	// fmt.Println(`----`)
	var expChAprov = regexp.MustCompile(`(?-s).*channelapproved.*(?s)`)
	chAprov := expChAprov.FindString(logString)
	// fmt.Println(chAprov)
	chAprovSf := strings.Fields(chAprov)
	if len(chAprovSf) > 0 {
		// fmt.Printf("length of chAprovSf is greater: %d\n", len(chAprovSf))

		baseRel := strings.Split(chAprovSf[2], "/")
		rel := strings.ReplaceAll(baseRel[0], "(", "")
		base := strings.ReplaceAll(baseRel[1], ")", "")

		aprStatus := strings.Split(chAprovSf[5], ".")
		// fmt.Println(aprStatus[1])
		// fmt.Printf("Channel with ID approved with status %s\n", aprStatus[1])

		state1 := SwapStatus{
			Status: aprStatus[1],
			State:  chAprovSf[3],
			SwapID: chAprovSf[0],
			Base:   base,
			Rel:    rel,
		}

		// fmt.Println("state 1:", state1)
		state1JSON, _ := json.Marshal(state1)
		// fmt.Println("Channel Approved")
		// fmt.Println("state1 JSON:", string(state1JSON))
		if answer == "single" {
			return string(state1JSON), nil
		}
		if answer == "full" {
			statuses = append(statuses, state1)
		}
	} else {
		// fmt.Printf("length of chAprovSf is lower: %d\n", len(chAprovSf))
	}

	// if len(aprovIDSf) > 0 {
	// 	// fmt.Printf("length of aprovIDSf is greater: %d\n", len(aprovIDSf))
	// } else {
	// 	fmt.Printf("length of aprovIDSf is lower: %d\n", len(aprovIDSf))
	// }

	// fmt.Println(`----`)
	var expAprovID = regexp.MustCompile(`(?m)approvalid.+$`)
	aprovID := expAprovID.FindString(logString)
	// fmt.Println(aprovID)
	aprovIDSf := strings.Fields(aprovID)
	if len(aprovIDSf) > 0 {
		// fmt.Printf("length of aprovIDSf is greater: %d\n", len(aprovIDSf))

		// fmt.Println(aprovIDSf)
		// fmt.Println(aprovIDSf[0])
		aprovStatus := strings.Split(aprovIDSf[0], ".")
		// fmt.Println(aprovStatus[0])
		// fmt.Println(aprovStatus[1])
		// fmt.Println(aprovIDSf[1])

		state1 := SwapStatus{
			State:     aprovStatus[0],
			Status:    "1",
			StateHash: aprovIDSf[1],
		}

		// fmt.Println("state 1:", state1)
		state1JSON, _ := json.Marshal(state1)
		// fmt.Println("Channel Approval ID")
		// fmt.Println("state1 JSON:", string(state1JSON))
		if answer == "single" {
			return string(state1JSON), nil
		}
		if answer == "full" {
			statuses = append(statuses, state1)
		}
	} //else {
	// fmt.Printf("length of aprovIDSf is lower: %d\n", len(aprovIDSf))
	//}

	// fmt.Println(`----`)
	var expIncCh = regexp.MustCompile(`(?m)incomingchannel.+$`)
	incCh := expIncCh.FindString(logString)
	// fmt.Println(incCh)
	incChSf := strings.Fields(incCh)
	// fmt.Println(incChSf)
	if len(incChSf) > 0 {
		// fmt.Printf("length of incChSf is greater: %d\n", len(incChSf))

		// fmt.Println(incChSf[0])
		incChStatus := strings.Split(incChSf[1], ".")
		// fmt.Println(incChStatus[0])
		// fmt.Println(incChStatus[1])

		state2 := SwapStatus{
			State:  incChSf[0],
			Status: incChStatus[1],
		}

		// fmt.Println("state 1:", state2)
		state2JSON, _ := json.Marshal(state2)
		// fmt.Println("Incoming Channel")
		// fmt.Println("state2 JSON:", string(state2JSON))
		if answer == "single" {
			return string(state2JSON), nil
		}
		if answer == "full" {
			statuses = append(statuses, state2)
		}
	} else {
		// fmt.Printf("length of incChSf is lower: %d\n", len(incChSf))
	}

	// fmt.Println(`----`)
	var expGotTxID = regexp.MustCompile(`(?m)got txid.+$`)
	TxID := expGotTxID.FindString(logString)
	// fmt.Println(TxID)
	TxIDSf := strings.Fields(TxID)
	// fmt.Println(TxIDSf)
	if len(TxIDSf) > 0 {
		// fmt.Printf("length of TxIDSf is greater: %d\n", len(TxIDSf))

		// fmt.Println(TxIDSf[0])
		TxIDStatus := strings.Split(TxIDSf[1], ".")
		// fmt.Println(TxIDStatus[0])
		// fmt.Println(TxIDStatus[1])

		state3 := SwapStatus{
			State:    TxIDSf[0],
			BaseTxID: TxIDStatus[1],
			Status:   "3",
		}

		// fmt.Println("state 1:", state3)
		state3JSON, _ := json.Marshal(state3)
		// fmt.Println("Sending TxID")
		// fmt.Println("state3 JSON:", string(state3JSON))
		if answer == "single" {
			return string(state3JSON), nil
		}
		if answer == "full" {
			statuses = append(statuses, state3)
		}
	} //else {
	// fmt.Printf("length of TxIDSf is lower: %d\n", len(TxIDSf))
	//}

	// fmt.Println(`----`)
	var expZFrom = regexp.MustCompile(`(?m)from..+$`)
	zFrom := expZFrom.FindString(logString)
	// fmt.Println("\tzFrom:", zFrom)
	zFromSf := strings.Fields(zFrom)
	// fmt.Println("\tzFromSf:", zFromSf)

	if len(zFromSf) > 0 {
		fmt.Printf("length of zFromSf is greater: %d\n", len(zFromSf))

		fmt.Println(zFromSf[0])
		zFromSl := strings.Split(zFromSf[0], ".")
		zFromAddr := strings.ReplaceAll(zFromSl[1], "(", "")
		zFromAddr = strings.ReplaceAll(zFromAddr, ")", "")
		// fmt.Println("\tzFromAddr: ", zFromAddr)
		// fmt.Println("\tzFromSf[2]: ", zFromSf[2])
		zFromJSON := strings.ReplaceAll(zFromSf[2], "'", "")
		zFromJSON = strings.ReplaceAll(zFromJSON, "'", "")
		zFromJSON = strings.ReplaceAll(zFromJSON, "[\"", "[")
		zFromJSON = strings.ReplaceAll(zFromJSON, "]\"", "]")
		zFromJSON = strings.ReplaceAll(zFromJSON, "\\", "")
		// fmt.Printf("\tzFromJSON: %s\n", zFromJSON)
		var zj ZFrom
		err := json.Unmarshal([]byte(zFromJSON), &zj)
		if err != nil {
			log.Println(err)
			// log.Fatalln(err)
		}
		// fmt.Println("\tzj: ", zj)
		// fmt.Println("\tzj[0].Address: ", zj[0].Address)
		// fmt.Println("\tzj[0].Amount: ", zj[0].Amount)
		// fmt.Println("\tzj[0].Memo: ", zj[0].Memo)

		state3 := SwapStatus{
			State:      "Sending Z Transaction",
			Status:     "3",
			SendToAddr: zj[0].Address,
			BaseAmount: zj[0].Amount,
		}

		// fmt.Println("\tstate 1: ", state3)
		state3JSON, _ := json.Marshal(state3)
		// fmt.Println("\tSending Z tx")
		// fmt.Println("\tstate3 JSON: ", string(state3JSON))
		if answer == "single" {
			return string(state3JSON), nil
		}
		if answer == "full" {
			statuses = append(statuses, state3)
		}
	} //else {
	// fmt.Printf("length of zFromSf is lower: %d\n", len(zFromSf))
	//}

	// fmt.Println(`----`)
	var expIncPay = regexp.MustCompile(`(?m)incomingpayment.+$`)
	incPay := expIncPay.FindString(logString)
	// fmt.Println(incPay)
	incPaySf := strings.Fields(incPay)
	// fmt.Println(incPaySf)
	if len(incPaySf) > 0 {
		// fmt.Printf("length of incPaySf is greater: %d\n", len(incPaySf))

		// fmt.Println(incPaySf[0])
		incPayStatus := strings.Split(incPaySf[1], ".")
		// fmt.Println(incPayStatus[0])
		// fmt.Println(incPayStatus[1])

		state4 := SwapStatus{
			State:  incPaySf[0],
			Status: incPayStatus[1],
		}

		// fmt.Println("state 1:", state4)
		state4JSON, _ := json.Marshal(state4)
		// fmt.Println("Incoming Payment")
		// fmt.Println("state4 JSON:", string(state4JSON))
		if answer == "single" {
			return string(state4JSON), nil
		}
		if answer == "full" {
			statuses = append(statuses, state4)
		}
	} //else {
	// fmt.Printf("length of incPaySf is lower: %d\n", len(incPaySf))
	//}

	// fmt.Println(`----`)
	var expAliceWait = regexp.MustCompile(`(?m)alice waits.+$`)
	aliceWait := expAliceWait.FindString(logString)
	// fmt.Println(aliceWait)
	aliceWaitSf := strings.Fields(aliceWait)
	// fmt.Println(aliceWaitSf)
	if len(aliceWaitSf) > 0 {
		// fmt.Printf("length of aliceWaitSf is greater: %d\n", len(aliceWaitSf))

		// fmt.Println(aliceWaitSf[3])
		aliceWaitTxID := strings.Split(aliceWaitSf[3], ".")
		// fmt.Println(aliceWaitTxID[0])
		// fmt.Println(aliceWaitTxID[1])

		// fmt.Println(aliceWaitSf[8])
		rcvAmount := strings.ReplaceAll(aliceWaitSf[8], "(", "")
		// fmt.Println(rcvAmount)
		// fmt.Println(aliceWaitSf[10])
		rcvAddr := strings.ReplaceAll(aliceWaitSf[10], ")", "")
		// fmt.Println(rcvAddr)
		rcvAmountflt, _ := strconv.ParseFloat(rcvAmount, 64)

		state4 := SwapStatus{
			State:     "incomingpayment",
			Status:    "4",
			RelAmount: rcvAmountflt,
			RelTxID:   aliceWaitTxID[1],
			Rel:       aliceWaitTxID[0],
			RecvAddr:  rcvAddr,
		}

		// fmt.Println("state 1:", state4)
		state4JSON, _ := json.Marshal(state4)
		// fmt.Println("Alice Waiting Payment")
		// fmt.Println("state4 JSON:", string(state4JSON))
		if answer == "single" {
			return string(state4JSON), nil
		}
		if answer == "full" {
			statuses = append(statuses, state4)
		}
	} //else {
	// fmt.Printf("length of aliceWaitSf is lower: %d\n", len(aliceWaitSf))
	//}

	// fmt.Println(`----`)
	var expAliceRcvd = regexp.MustCompile(`(?m)received.+$`)
	aliceRcvd := expAliceRcvd.FindString(logString)
	// fmt.Println(aliceRcvd)
	aliceRcvdSf := strings.Fields(aliceRcvd)
	// fmt.Println(aliceRcvdSf)
	if len(aliceRcvdSf) > 0 {

		// fmt.Println(aliceRcvdSf[1])
		// fmt.Println(aliceRcvdSf[3])
		rcvAmountflt, _ := strconv.ParseFloat(aliceRcvdSf[1], 64)

		state4 := SwapStatus{
			State:     aliceRcvdSf[0],
			Status:    "4",
			RelAmount: rcvAmountflt,
		}

		// fmt.Println("state 1:", state4)
		state4JSON, _ := json.Marshal(state4)
		// fmt.Println("Alice Received Payment")
		// fmt.Println("state4 JSON:", string(state4JSON))
		if answer == "single" {
			return string(state4JSON), nil
		}
		if answer == "full" {
			statuses = append(statuses, state4)
		}
	} //else {
	// fmt.Printf("length of aliceRcvdSf is lower: %d\n", len(aliceRcvdSf))
	//}

	// fmt.Println(`----`)
	var expSwpCompl = regexp.MustCompile(`(?m)SWAP COMPLETE.+$`)
	swpCompl := expSwpCompl.FindString(logString)
	// fmt.Println(swpCompl)
	swpComplSf := strings.Fields(swpCompl)
	// fmt.Println(swpComplSf)
	if len(swpComplSf) > 0 {

		// fmt.Println(swpComplSf[0])
		// fmt.Println(swpComplSf[1])

		state4 := SwapStatus{
			State:  swpComplSf[0] + string(' ') + swpComplSf[1],
			Status: "4",
		}

		state4JSON, _ := json.Marshal(state4)
		// fmt.Println("SWAP COMPLETE")
		// fmt.Println("state4 JSON:", string(state4JSON))
		if answer == "single" {
			return string(state4JSON), nil
		}
		if answer == "full" {
			statuses = append(statuses, state4)
		}
	} //else {
	// fmt.Printf("length of swpComplSf is lower: %d\n", len(swpComplSf))
	//}

	// fmt.Println(`----`)
	var expIncPaid = regexp.MustCompile(`(?m)incomingfullypaid.+$`)
	incPaid := expIncPaid.FindString(logString)
	// fmt.Println(incPaid)
	incPaidSf := strings.Fields(incPaid)
	// fmt.Println(incPaidSf)
	if len(incPaidSf) > 0 {

		// fmt.Println(incPaidSf[0])
		// fmt.Println(incPaidSf[1])
		incPaidStatus := strings.Split(incPaidSf[1], ".")
		// fmt.Println(incPaidStatus[1])

		state5 := SwapStatus{
			State:  incPaidSf[0],
			Status: incPaidStatus[1],
		}

		state5JSON, _ := json.Marshal(state5)
		// fmt.Println("SWAP COMPLETE")
		// fmt.Println("state5 JSON:", string(state5JSON))
		if answer == "single" {
			return string(state5JSON), nil
		}
		if answer == "full" {
			statuses = append(statuses, state5)
		}
	} //else {
	// fmt.Printf("length of incPaidSf is lower: %d\n", len(incPaidSf))
	//}

	// fmt.Println(`----`)
	var expIncClose = regexp.MustCompile(`(?m)incomingclose.+$`)
	incClose := expIncClose.FindString(logString)
	// fmt.Println(incClose)
	incCloseSf := strings.Fields(incClose)
	// fmt.Println(incCloseSf)
	if len(incCloseSf) > 0 {

		// fmt.Println(incCloseSf[0])
		// fmt.Println(incCloseSf[1])
		incCloseStatus := strings.Split(incCloseSf[1], ".")
		// fmt.Println(incCloseStatus[1])

		state6 := SwapStatus{
			State:  incCloseSf[0],
			Status: incCloseStatus[1],
		}

		state6JSON, _ := json.Marshal(state6)
		// fmt.Println("SWAP COMPLETE")
		// fmt.Println("state6 JSON:", string(state6JSON))
		if answer == "single" {
			return string(state6JSON), nil
		}
		if answer == "full" {
			statuses = append(statuses, state6)
		}
	} //else {
	// fmt.Printf("length of incCloseSf is lower: %d\n", len(incCloseSf))
	//}

	// fmt.Println(`----`)
	var expOpid = regexp.MustCompile(`(?m)opid..+$`)
	opID := expOpid.FindString(logString)
	// fmt.Println(opID)
	opIDSf := strings.Fields(opID)
	// fmt.Println(opIDSf[0])

	if len(opIDSf) > 0 {

		opIDSs := strings.Split(opIDSf[0], ".")
		// fmt.Println(opIDSs[1])
		opIDRa := strings.ReplaceAll(opIDSs[1], "(", "")
		opIDRa = strings.ReplaceAll(opIDRa, ")", "")
		// fmt.Println(opIDRa)

		state6 := SwapStatus{
			State:    "opid",
			Status:   "6",
			BaseTxID: opIDRa,
		}

		state6JSON, _ := json.Marshal(state6)
		// fmt.Println("SWAP COMPLETE")
		// fmt.Println("state6 JSON:", string(state6JSON))
		if answer == "single" {
			return string(state6JSON), nil
		}
		if answer == "full" {
			statuses = append(statuses, state6)
		}
	} //else {
	// fmt.Printf("length of opIDSf is lower: %d\n", len(opIDSf))
	//}

	// fmt.Println(`----`)
	var expDpow = regexp.MustCompile(`(?m)dpow_broadcast.+$`)
	dPowBcast := expDpow.FindString(logString)
	// fmt.Println(dPowBcast)
	dPowBcastSf := strings.Fields(dPowBcast)
	// fmt.Println(dPowBcastSf)
	if len(dPowBcastSf) > 0 {

		// fmt.Println(dPowBcastSf[0])
		// fmt.Println(dPowBcastSf[3])

		state6 := SwapStatus{
			State:    "dpow_broadcast",
			Status:   "6",
			BaseTxID: dPowBcastSf[3],
		}

		state6JSON, _ := json.Marshal(state6)
		// fmt.Println("SWAP COMPLETE")
		// fmt.Println("state6 JSON:", string(state6JSON))
		if answer == "single" {
			return string(state6JSON), nil
		}
		if answer == "full" {
			statuses = append(statuses, state6)
		}
	} //else {
	// fmt.Printf("length of dPowBcastSf is lower: %d\n", len(dPowBcastSf))
	//}

	// fmt.Println(`----`)
	var expSwapData = regexp.MustCompile(`(?-s).*subatomic_cmd.*(?s)`)
	SwapData := expSwapData.FindString(logString)
	// fmt.Println(SwapData)
	stateDone := SwapStatus{
		SwapFullData: SwapData,
	}

	stateDoneJSON, _ := json.Marshal(stateDone)
	if answer == "single" {
		return string(stateDoneJSON), nil
	}
	if answer == "full" {
		statuses = append(statuses, stateDone)
	}

	statusesJSON, _ := json.Marshal(statuses)
	// fmt.Println(statusesJSON)

	if answer == "full" {
		return string(statusesJSON), nil
	}

	return "{}", errors.New("no log found")

}
