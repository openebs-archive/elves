import os,time

i = 0
while i < 1000:
    time.sleep(10)
    os.system("python generate_io.py")
    i += 1
    os.command("sync;sync;sync")
    time.sleep(20)
    os.remove("file.txt")

