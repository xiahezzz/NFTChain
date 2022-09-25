networkup:
	./network.sh up

networkdown:
	./network.sh down

createChannel:
	./network.sh createChannel -c channel

deployCC:
	./network.sh deployCC -ccn test2 -ccp ../chaincode -ccl go -c channel

.PHONY:networkup networkdown createChannel deployCC