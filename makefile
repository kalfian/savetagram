build:
	go build -o build/savetagram	

serve-god:
	pkill -9 savetagram
	rm savetagram.log
	go build -o build/savetagram
	god --nohup --logfile savetagram.log --rundir /home/kalfianc/savetagram -- ./build/savetagram

.PHONY: build serve-god