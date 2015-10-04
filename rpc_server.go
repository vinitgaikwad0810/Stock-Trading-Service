	// server.go
	package main
	import (
"log"
"net"
"net/rpc"
"net/rpc/jsonrpc"
	//"log"
"net/http"
//"os"
"strings"
"fmt"
"strconv"
"net/url"
"math"
ejson "encoding/json"
)
	var TRACK_ID_CONSTANT int
	var response_bought_stocks [100]Response_bought_stocks 	
	type Args struct {
		X, Y int
	}
	type Buying_stocks_input struct {
		StockSymbolAndPercentage string
		Budget float64
	}

	type Response_bought_stocks struct {
		trackId int
		stocks []string
		unvestedAmount float64
	}
	type Stocks struct{
		company_name string
		number_of_shares int
		stock_price string
	}
	type Response_output struct {
		Message string
	}
	type Stocktradingsystem struct {}


	func BuySharesRealTime (Budget float64, company string) (string,float64){
		var Url *url.URL
		Url, err := url.Parse("https://query.yahooapis.com")
		if err != nil {
			panic("Error Panic")
		}
		Url.Path += "/v1/public/yql"
		parameters := url.Values{}
		parameters.Add("q", "select * from yahoo.finance.quote where symbol in ('"+company+"')")
		Url.RawQuery = parameters.Encode()
		Url.RawQuery += "&format=json&diagnostics=true"
		Url.RawQuery += "&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="
	//fmt.Printf("Encoded URL is %q\n", Url.String())
		res, err := http.Get(Url.String())
									//
		if err != nil {
			panic("Error Panic")
		}
		defer res.Body.Close()
		var v map[string]interface{}
		dec:= ejson.NewDecoder(res.Body);
		if err := dec.Decode(&v); err != nil {
			fmt.Println("ERROR: " + err.Error())
		}   
	   // person := new(Person)
	   // json.Unmarshal([]byte(res.String()), person)
		symbol := v["query"].(map[string]interface{})["results"].(map[string]interface{})["quote"].(map[string]interface{})["symbol"]
		lastTradePriceOnly := v["query"].(map[string]interface{})["results"].(map[string]interface{})["quote"].(map[string]interface{})["LastTradePriceOnly"]
		fmt.Println("The symbol is ",symbol,"LastTradePriceOnly =",lastTradePriceOnly)
		f, err := strconv.ParseFloat(lastTradePriceOnly.(string), 64)
		number_of_shares :=  int(Budget/f)
	//number_of_shares_float:= float64(number_of_shares)
		unvestedAmount := Budget - float64(number_of_shares)*f
		share := company +":"+strconv.Itoa(number_of_shares)+":$"+lastTradePriceOnly.(string)
		//fmt.Println("Share String is "+ share)
		return share,unvestedAmount
	}


	func (t *Stocktradingsystem) CheckingPortfolio(trackId int, reply *Response_output) error {

		fmt.Println("\n--------------PORTFOLIO LOGS----------------------------\n\n")
	//fmt.Println("Track ID =",trackId)
		reply.Message = "\nstocks -->"
		var currentMarketValue  float64
		fmt.Println(response_bought_stocks[trackId])
		for _,each_company := range response_bought_stocks[trackId].stocks{
		//fmt.Println(each_company)
			values := strings.Split(each_company,":")
		//fmt.Println(values[2])
			company := values[0]
			number_of_shares := values[1]
			stock_value,err := strconv.ParseFloat(strings.Trim(values[2], "$"),64)
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Println(stock_value)
			var Url *url.URL
			Url, err = url.Parse("https://query.yahooapis.com")
			if err != nil {
				panic("Error Panic")
			}
			Url.Path += "/v1/public/yql"
			parameters := url.Values{}
			 parameters.Add("q", "select * from yahoo.finance.quote where symbol in ('"+company+"')")//
			 Url.RawQuery = parameters.Encode()
			 Url.RawQuery += "&format=json&diagnostics=true"
			 Url.RawQuery += "&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="
			 res, err := http.Get(Url.String())
									//
			 if err != nil {
			 	panic("Error Panic")
			 }
			 defer res.Body.Close()
			 var v map[string]interface{}
			 dec:= ejson.NewDecoder(res.Body);
			 if err := dec.Decode(&v); err != nil {
			 	fmt.Println("ERROR: " + err.Error())
			 }   
			 symbol := v["query"].(map[string]interface{})["results"].(map[string]interface{})["quote"].(map[string]interface{})["symbol"]
			 lastTradePriceOnly := v["query"].(map[string]interface{})["results"].(map[string]interface{})["quote"].(map[string]interface{})["LastTradePriceOnly"]
			 fmt.Println("The symbol is ",symbol,"LastTradePriceOnly =",lastTradePriceOnly)
			 float_TradePrice, err := strconv.ParseFloat(lastTradePriceOnly.(string), 64)
			 if err != nil {
			 	panic("Error Panic")
			 }
			 var ReplyString string
			 	//float_TradePrice = float_TradePrice - 1.0
			 fmt.Println("stock_value =",stock_value," float_TradePrice =",float_TradePrice)
			 if stock_value < float_TradePrice{
			 	ReplyString = company +":"+ number_of_shares +":"+ "+$"+lastTradePriceOnly.(string) + " "
			 }else if stock_value > float_TradePrice{
			 	ReplyString = company +":"+ number_of_shares +":"+"-$"+lastTradePriceOnly.(string) + " "
			 }else if stock_value == float_TradePrice{
			 	ReplyString = company +":"+ number_of_shares +":"+"=$"+lastTradePriceOnly.(string)+" "
			 }
			 reply.Message += ReplyString
			 //MarketValue
			 currentMarketValue_this,err := strconv.ParseFloat(number_of_shares,64)
			 if err != nil {
			 	panic("Error Panic")
			 }
			 currentMarketValue += currentMarketValue_this*float_TradePrice
			}
			reply.Message += "\ncurrentMarketValue -->"+strconv.FormatFloat(currentMarketValue,'f',6,64)
			reply.Message += "\nunvestedAmount -->"+strconv.FormatFloat(response_bought_stocks[trackId].unvestedAmount,'f',6,64)

			fmt.Println("\n--------------PORTFOLIO LOGS END----------------------------\n\n")
			return nil
		}


		func (t *Stocktradingsystem) Buyingstocks(args *Buying_stocks_input, reply *Response_output) error {

			fmt.Println("\n--------------BUYING STOCKS LOGS----------------------------\n\n")
			result := strings.Split(args.StockSymbolAndPercentage, ",")
			response_bought_stocks[TRACK_ID_CONSTANT].trackId = TRACK_ID_CONSTANT
			for _,company := range result {
				//fmt.Println("Index is ",index)
				stock_data := strings.Split(company,":");
				//fmt.Println(stock_data)
				company_name := stock_data[0]
				percentages := stock_data[1]
				allocated_budget,err := ParseFloatPercent(percentages,64) 
				if err != nil {
					fmt.Println(err)
				}
				allocated_budget_amount := allocated_budget * args.Budget
				fmt.Println("allocated_budget = " + FloatToString(allocated_budget))
				fmt.Println("allocated_budget_amount = " + FloatToString(allocated_budget_amount))
				number_of_shares,unvestedAmount := BuySharesRealTime(allocated_budget_amount, company_name)
				response_bought_stocks[TRACK_ID_CONSTANT].unvestedAmount += unvestedAmount
				fmt.Println("number_of_shares = ",number_of_shares ,"unvestedAmount= ",response_bought_stocks[TRACK_ID_CONSTANT].unvestedAmount)
				response_bought_stocks[TRACK_ID_CONSTANT].stocks = append(response_bought_stocks[TRACK_ID_CONSTANT].stocks,number_of_shares)
			}
			stocks:=serialize(response_bought_stocks[TRACK_ID_CONSTANT].stocks)
			reply.Message = "\ntrackId ="+ strconv.Itoa(response_bought_stocks[TRACK_ID_CONSTANT].trackId) + "\nstocks =" + stocks + "\nunvestedAmount =" + strconv.FormatFloat(response_bought_stocks[TRACK_ID_CONSTANT].unvestedAmount,'f',6,64)
			TRACK_ID_CONSTANT = TRACK_ID_CONSTANT + 1
			//fmt.Println(response_bought_stocks[TRACK_ID_CONSTANT])
			fmt.Println("\n--------------BUYING STOCKS LOGS ENDED----------------------------\n\n")
			return nil
		}


		func serialize(input []string) string{
			var output string
			//length := len(input)
			for _,stringval := range input{
				stringval = " " + stringval
				output += stringval
			}
			//fmt.Println("output" + output)
			return output
		}
		func FloatToString(input_num float64) string {
	    // to convert a float number to a string
			return strconv.FormatFloat(input_num, 'f', 6, 64)
		}


		func ParseFloatPercent(s string, bitSize int) (f float64, err error) {
			i := strings.Index(s, "%")
			if i < 0 {
				return 0, fmt.Errorf("ParseFloatPercent: percentage sign not found")
			}
			f, err = strconv.ParseFloat(s[:i], bitSize)
			if err != nil {
				return 0, err
			}
			return f / 100, nil
		}


		func RoundPlus(f float64, places int) (float64) {
			shift := math.Pow(10, float64(places))
			return Round(f * shift) / shift;    
		}

		func Round(f float64) float64 {
			return math.Floor(f + .5)
		}

		func main() {

			fmt.Println("\n--------------SERVER LOGS ARE ENABLED----------------------------\n\n")
			TRACK_ID_CONSTANT = 1
			cal := new(Stocktradingsystem)
			server := rpc.NewServer()
			server.Register(cal)
			server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
			listener, e := net.Listen("tcp", ":1234")
			if e != nil {
				log.Fatal("listen error:", e)
			}
			for {
				if conn, err := listener.Accept(); err != nil {
					log.Fatal("accept error: " + err.Error())
				} else {
					log.Printf("new connection established\n")
					go server.ServeCodec(jsonrpc.NewServerCodec(conn))
				}
			}
		}
