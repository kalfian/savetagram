build:
	go build -o savetagram

serve-god:
	god --nohup --logfile savetagram.log --rundir /home/kalfianc/savetagram -- ./savetagram