# cmpe273-assignment1
Stock Trading Server using JSON RPC and golang





1. Buying stocks: Terminal or console input and output of the Client.go.
Your Answer:

vinit@ubuntu:~/golang/jsonrpcfinal$ go run rpc_client.go GOOG:50%,YHOO:50%  20000

You have opted to buy stocks (PART 1)

trackId =5
stocks = GOOG:15:$626.91 YHOO:325:$30.71
unvestedAmount =615.600000
vinit@ubuntu:~/golang/jsonrpcfinal$




2. Checking your portfolio: Terminal or console input and output of the Client.go.
Your Answer:

vinit@ubuntu:~/golang/jsonrpcfinal$ go run rpc_client.go 5
You have opted for portfolio (PART 2)


stocks -->GOOG:15:=$626.91 YHOO:325:=$30.71
currentMarketValue -->19384.400000
unvestedAmount -->615.600000
vinit@ubuntu:~/golang/jsonrpcfinal$

