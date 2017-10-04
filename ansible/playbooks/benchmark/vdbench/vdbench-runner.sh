#!/bin/bash

########################################################################################
# This script runs vdbench I/O in the vdbench test container, on /datadir mounted volume 
# for the workload configs placed in /templates and plots the IOPS, BW, Lat & CPU 
# utilization graphs. The logs and plots are placed in /logdir
#
# The container has to be launched with volumes (host dirs) holding templates, datadir 
# and logdir 
########################################################################################

TEMPLATE_DIR="/templates"
LOG_DIR="/logdir"
DATA_DIR="/datadir"

if [ $# -gt 0 ]; 
then
  echo "Setting the templates directory as $1"
  TEMPLATE_DIR=$1
  if [ ! -d "$TEMPLATE_DIR" ]; 
  then
    echo "Specified template dir doesn't exist"
    exit 1 
  fi
fi 

if [ $# -gt 1 ];
then
  echo "Settings the log directory as $2"
  LOG_DIR=$2
  if [ ! -d "$LOG_DIR" ];
  then
    echo "Specified log dir doesn't exist"
    exit 1
  fi
fi 

# Verify that the data directory is mounted

df -h -P | grep -qw "datadir"
if [ $? -ne 0 ]; 
then
  echo -e "datadir is not mounted in container, exiting \n"
  exit 
else
  echo "datadir mounted successfully"
fi


for i in `ls ${TEMPLATE_DIR}`
do
  TEST_NAME=`echo $i | cut -d "." -f 1`
  SAMPLES=`grep elapsed ${TEMPLATE_DIR}/$i | cut -d "," -f 3 | cut -d "=" -f 2`
  OUT_DIR="${LOG_DIR}/${TEST_NAME}"
  
  mkdir -p $OUT_DIR
  if [ $? -ne 0 ]; 
  then 
    echo -e "unable to create $OUT_DIR, exiting \n"
    exit
  fi
 
  echo -e "Running $i workload\n"  

  ./vdb/vdbench -f ${TEMPLATE_DIR}/$i -o $OUT_DIR
  if [ $? -ne 0 ]; then
    echo "vdbench failed to run, exiting"  
    exit
  fi

  ./vdb2gnuplot.sh -f ${OUT_DIR}/flatfile.html -n $SAMPLES -a 
  if [ $? -ne 0 ]; then
    echo "Failed to generate graph, exiting"
    exit
  fi 

  mv *.png $OUT_DIR/
done
 



