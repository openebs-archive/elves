from __future__ import division
import subprocess
import time, os

cmd_total_mem = "sudo free | awk ' /'Mem'/ {print $2}'"

out = subprocess.Popen(cmd_total_mem,stdout=subprocess.PIPE, shell=True)
total_mem = out.communicate()
list = []
n = 10
count = 0
while count < n:
    count = count + 1
    cmd_used_mem = "sudo free | awk ' /'Mem'/ {print $3}'"
    out = subprocess.Popen(cmd_used_mem,stdout=subprocess.PIPE, shell=True)
    used_mem = out.communicate()
    mem_utilized = int(total_mem[0])/int(used_mem[0])
    time.sleep(20)
    list.append(mem_utilized)
if all(i <= 30 for i in list):
        print "Test Passed"
else:
        print "Test Failed"

