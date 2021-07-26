build:
	go build -o build/savetagram

serve-god:
	god --nohup --logfile savetagram.log --rundir /home/kalfianc/savetagram -- ./build/savetagram