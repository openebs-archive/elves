#!/bin/bash

# CDDL HEADER START
#
# This file and its contents are supplied under the terms of the
# Common Development and Distribution License ("CDDL"), version 1.0.
# You may only use this file in accordance with the terms of version
# 1.0 of the CDDL.
#
# A full copy of the text of the CDDL should have accompanied this
# source.  A copy of the CDDL is also available via the Internet at
# http://www.illumos.org/license/CDDL.
#
# CDDL HEADER END
#
# Copyright 2015 Nexenta Systems Inc.
#
# wget --no-ch https://bitbucket.org/alek_p/vdb2gnuplot/raw/HEAD/vdb2gnuplot.sh

VDB2GNUPLOT_VERSION="1.5"

VDBENCH="./vdb/vdbench"

# exits script
usage() {
	echo "$0 v$VDB2GNUPLOT_VERSION Usage:\n"
	echo "$0 -f path/to/flatfile.html -n #.of_samples_in_run [-i|-b|-l|-c|-a] [-R #,#,#...] [-d str.data_displ] [-p str.png_prefix] [-r #.plot_this_rdpct] [-t #.plot_this_threads] [-s #.bytes_plot_this_xfersize] [-k]"
	echo "\t -f path to flatflie.html (required)"
	echo "\t -n samples per run ( elapsed / interval ) (required)"
	echo "\t -i plot IOPS graph (default)"
	echo "\t -b plot MB/s graph"
	echo "\t -l plot latency (ms) graph"
	echo "\t -c plot cpu utilization (usr+sys) % graph"
	echo "\t -R plot only specified runs"
	echo "\t -d gnuplot data display format (default is linespoints)"
	echo "\t -p optional prefix for png file"
	echo "\t -r only plot runs that have specified \"rdpct\" used (supported for Storage definitions only)"
	echo "\t -t only plot runs that have specified \"threads\" used (supported for Storage definitions only)"
	echo "\t -s only plot runs that have specified \"xfersize\" used"
	echo "\t -k put the graph key above title"
	echo "\t -a plot all graphs (-l -i -b -c)"

	exit 1
}

FLATFILE_FIRST_DATA=1
flatfile_first_data() {
	while read line
	do
		if [[ "$line" =~ ^[0-9]+:[0-9]+:[0-9]+ ]]; then
			if [[ "$line" != *"format_for_"* ]]; then
				FLATFILE_FIRST_DATA=$((FLATFILE_FIRST_DATA+1))
				return
			fi
		fi
		FLATFILE_FIRST_DATA=$((FLATFILE_FIRST_DATA+1))
	done < $1

	usage
}

run_not_needed() {
	IFS=',' read -ra RUNS_WANTED <<< "$1"
	for J in "${RUNS_WANTED[@]}"; do
		if [ $J -eq "$2" ]; then
			return 0
		fi
	done

	return 1
}

timestamp_get_sec() {
	read -r h m s <<< $(echo 10#$1 | tr ':' ' ' )
	echo $(((${h#0}*60*60)+(${m#0}*60)+${s#0}))
}

check_util_in_path() {
	UTIL_PATH=$(which $1)
	if ! [ -x "$UTIL_PATH" ]
	then
		echo "$1 is a prereq and it was not found, exiting."
		usage
	fi
}

# default plot is latency
PLOT_TYPE=5
DATA_DISP_TYPE="linespoints"
OPTIND=1
ARGS=$*
while getopts "h?iblcp:R:r:t:s:f:n:d:ka" opt; do
	case "$opt" in
		i)  PLOT_TYPE=3 # IOPS RATE
		;;
		b)  PLOT_TYPE=4 # BANDWIDTH MB/S
		;;
		l)  PLOT_TYPE=5 # LATENCY MS
		;;
		c)  PLOT_TYPE=6 # CPU UTILIZATION %
		;;
		d)  DATA_DISP_TYPE=${OPTARG}
		;;
		p)  PLOT_FILENAME=${OPTARG}
		;;
		R)  RUN_WANTED=${OPTARG}
		;;
		s)  XS_WANTED=${OPTARG}
		;;
		r)  RDPCT_WANTED=${OPTARG}
		;;
		t)  THREADS_WANTED=${OPTARG}
		;;
		f)  FLATFILE=${OPTARG}
		;;
		n)  SAMPLES_IN_RUN=${OPTARG}
		;;
		k)  LEGEND_TOP=1
		;;
		a)  PLOT_TYPE=1 # ALL
		;;
		h|\?|*)
		    usage
		;;
	esac
done

# make sure prereqs are there
# check_util_in_path "vdbench"
# check_util_in_path "gnuplot"

shift $((OPTIND-1))
if [ -n "$1" ]; then
	echo "Extra cmd line args ($1) detected."
	usage
fi

if [ -z "${FLATFILE}" ] || [ -z "${SAMPLES_IN_RUN}" ]; then
	echo "Flatfile and/or Samples in Run not specified."
	usage
fi

if ! [ -f "${FLATFILE}" ]; then
	echo "Provided flatfile ($FLATFILE) does not exist."
	usage
fi

if [[ PLOT_TYPE -eq 1 ]]; then
	ARGS=`echo "$ARGS" | sed 's/\( -a\)//g'`
	$0 $ARGS -l
	$0 $ARGS -b
	$0 $ARGS -i
	$0 $ARGS -c

	exit $?
fi

flatfile_first_data $FLATFILE

# need to figure out seconds per run
FIRST_TIMESTAMP=`sed -n "$((FLATFILE_FIRST_DATA)),$((FLATFILE_FIRST_DATA))p" $FLATFILE | cut -f 1 -d "."`
SECOND_TIMESTAMP=`sed -n "$((FLATFILE_FIRST_DATA + 1)),$((FLATFILE_FIRST_DATA + 1))p" $FLATFILE | cut -f 1 -d "."`
FIRST_TIMESTAMP=$(timestamp_get_sec $FIRST_TIMESTAMP)
SECOND_TIMESTAMP=$(timestamp_get_sec $SECOND_TIMESTAMP)
SECONDS_PER_RUN=$((SECOND_TIMESTAMP-FIRST_TIMESTAMP))


FLATFILE_LINE25=`sed -n '25,25p' $FLATFILE`
if [ "$FLATFILE_LINE25" == "* rdpct           : read% requested" ]; then
	RUN_TYPE="Storage";
	READ_STAT_TYPE="rqread"
else
	RUN_TYPE="FS"
	READ_STAT_TYPE="ks_read%"
fi

echo -n "\n$RUN_TYPE Definition was used. Selected plot type is "
if (($PLOT_TYPE==3)); then
	YLABEL="IOPS"
elif (($PLOT_TYPE==4)); then
	YLABEL="MB_per_s"
elif (($PLOT_TYPE==5)); then
	YLABEL="latency_in_ms"
#	YRANGE="set yrange [0:15]"
elif (($PLOT_TYPE==6)); then
	YLABEL="CPU_pcnt_used"
fi
echo $YLABEL
echo " "

echo "Seconds per run: $SECONDS_PER_RUN"

CSV_TMP="tmp.vdb2gnuplot.$$.flatfile.csv"

#PARSEFLAT_STR="$VDBENCH parseflat -i $FLATFILE -o $CSV_TMP -c run -c interval -c rate -c MB/sec -c resp -c cpu_used -c xfersize -c $READ_STAT_TYPE"

PARSEFLAT_STR="$VDBENCH parseflat -i $FLATFILE -o $CSV_TMP -c run -c interval -c rate -c MB/sec -c resp -c cpu_used -c xfersize"

if [ -n "${XS_WANTED}" ]; then
	PLOT_FILENAME+="xfersize_$XS_WANTED."
fi

if [ -n "${RDPCT_WANTED}" ]; then
	if [ $RUN_TYPE == "FS"  ]; then
		echo "'-r #' is not supported for FS definition as rdpct is not reported in the flatfile, skipping it."
	else
		PLOT_FILENAME+="rdpct_$RDPCT_WANTED."
	fi
fi

if [ -n "${THREADS_WANTED}" ]; then
	if [ $RUN_TYPE == "FS"  ]; then
		echo "'-t #' is not supported for FS definition as threads is not reported in the flatfile, skipping it."
	else
		PARSEFLAT_STR+=" -c threads"
		PLOT_FILENAME+="threads_$THREADS_WANTED."
	fi
fi

PLOT_FILENAME+="$YLABEL.png"
PARSEFLAT_STR+=" 1>/dev/null 2>&1"

echo -n "Generating CSV file ($CSV_TMP) from $FLATFILE... "
eval $PARSEFLAT_STR
if [ $? -ne 0 ]; then
        echo "FAILED."
	`$PARSEFLAT_STR`
        exit 1
else
        echo "done."
fi

# drop header
perl -ni -e 'print unless $. == 1' $CSV_TMP

TOTAL_LINES=`wc -l < "$CSV_TMP" | tr -d ' '`

NUM_RUNS=$((TOTAL_LINES/SAMPLES_IN_RUN))

echo "Total data rows "$TOTAL_LINES" with "$SAMPLES_IN_RUN" rows per run we have "$NUM_RUNS" runs."

if [[ $NUM_RUNS -eq 0 ]]; then
	echo "Samples per run is more than the number of actuall data points. Exiting."
	echo "\nNOTE:\nVDbench may have droped some data. Check output for:"
	echo "\"Detailed reporting is running behind; reporting of intervals #-# has been skipped.\""
	echo "You may be able to adjust samples per run if you still want to plot the data you have. (-n $TOTAL_LINES)."
	rm -f $CSV_TMP
	exit 1
fi

FLATFILE_PREFIX=`basename $FLATFILE`
if [[ $FLATFILE_PREFIX == *"flatfile.html"* ]]; then
	PLOT_FILENAME="$PLOT_FILENAME"
else
	FLATFILE_PREFIX=`echo $FLATFILE_PREFIX | sed 's/.html//'`
	PLOT_FILENAME="$FLATFILE_PREFIX.$PLOT_FILENAME"
fi

GNUPLOT_TMP="tmp.vdb2gnuplot.$$.$PLOT_FILENAME.gnuplot_cmds"
echo "Generating $YLABEL gnuplot command file ($GNUPLOT_TMP)"
echo "reset
set datafile separator \",\"
set title \"$YLABEL plot ($PLOT_FILENAME)\"
set output \"$PLOT_FILENAME\"
set xlabel \"Seconds\"
set ylabel \"$YLABEL\"
$YRANGE
set terminal png size 1024,768
set grid nopolar
set grid layerdefault linetype 0 linewidth 1, linetype 0 linewidth 1" > $GNUPLOT_TMP

if [ -n "${LEGEND_TOP}" ]; then
	echo "set key outside" >> $GNUPLOT_TMP
	echo "set key center top" >> $GNUPLOT_TMP
fi


echo "plot \\" >> $GNUPLOT_TMP

SED_RUN_START=1
for ((I=1; I<=$NUM_RUNS; I++, SED_RUN_START=$((SED_RUN_START+SAMPLES_IN_RUN))))
do
	echo -n "Processing $I / $NUM_RUNS runs."
	RUN_LABEL=`sed -n $SED_RUN_START,"${SED_RUN_START}p" $CSV_TMP | cut -d, -f1` # run name
	# process -R
	if [ -n "${RUN_WANTED}" ]; then
		run_not_needed "$RUN_WANTED" "$I"
		if (($? == 1)); then
			echo " Skipping run $I - $RUN_LABEL."
			continue
		else
			echo -n "\n"
		fi
	else
		echo -n "\n"

	fi
	# process -s
	if [ -n "${XS_WANTED}" ]; then
		XS_THIS_RUN=`sed -n $SED_RUN_START,"${SED_RUN_START}p" $CSV_TMP | cut -d, -f7`
		if ((XS_WANTED != XS_THIS_RUN)); then
			continue
		fi
	fi
	# process -r
	if [ -n "${RDPCT_WANTED}" ]; then
		RDPCT_THIS_RUN=`sed -n $SED_RUN_START,"${SED_RUN_START}p" $CSV_TMP | cut -d, -f8`
		RDPCT_THIS_RUN=${RDPCT_THIS_RUN%.*}
		RUN_LABEL+=" --- rdpct: $RDPCT_THIS_RUN; "
		if ((RDPCT_WANTED != RDPCT_THIS_RUN)); then
			continue
		fi
	fi
	# process -t
	if [ -n "${THREADS_WANTED}" ]; then
		THREADS_THIS_RUN=`sed -n $SED_RUN_START,"${SED_RUN_START}p" $CSV_TMP | cut -d, -f9`
		THREADS_THIS_RUN=${THREADS_THIS_RUN%.*}
		RUN_LABEL+=" --- threads: $THREADS_THIS_RUN;"
		if ((THREADS_WANTED != THREADS_THIS_RUN)); then
			continue
		fi
	fi

	SED_RUN_END="$((SAMPLES_IN_RUN*I))p"

	# sum of all xfersizes
	AVG=`sed -n $SED_RUN_START,$SED_RUN_END $CSV_TMP | cut -d, -f7 | paste -s -d+ | perl -nle 'print eval'`
	# avg xfersizes
	AVG=`echo "scale=1; $AVG / $SAMPLES_IN_RUN" | bc`
	#AVG=${AVG%.*}
	RUN_LABEL+=" --- xfersize: $AVG"

	# sum of all read %s
	AVG=`sed -n $SED_RUN_START,$SED_RUN_END $CSV_TMP | cut -d, -f8 | paste -s -d+ | perl -nle 'print eval'`
	# avg for read %
	AVG=`echo "scale=1; $AVG / $SAMPLES_IN_RUN" | bc`
	RUN_LABEL+=" --- read%: $AVG"

	# sum of all values for this plot type
	AVG=`sed -n $SED_RUN_START,$SED_RUN_END $CSV_TMP | cut -d, -f$PLOT_TYPE | paste -s -d+ | perl -nle 'print eval'`
	# avg for the whole run this plot type
	AVG=`echo "scale=2; $AVG / $SAMPLES_IN_RUN" | bc`

	PLOT_STR="\"<(sed -n $SED_RUN_START,$SED_RUN_END $CSV_TMP)\" using (\$2*$SECONDS_PER_RUN):$PLOT_TYPE title '$I: [$RUN_LABEL] $YLABEL avg: $AVG' with $DATA_DISP_TYPE lt $I lw 2, \\"
	echo $PLOT_STR >> $GNUPLOT_TMP
done

# remove the last ", \"
sed -i '$s/, \\//' $GNUPLOT_TMP

# plot the generated tmp file
echo "Generating $YLABEL gnuplot plot ($PLOT_FILENAME)"
gnuplot $GNUPLOT_TMP

# remove tmp files
rm -f $CSV_TMP
if (($? == 0)); then
	echo "Removed tmp file $CSV_TMP"
fi

rm -f $GNUPLOT_TMP
if (($? == 0)); then
	echo "Removed tmp file $GNUPLOT_TMP"
fi

exit 0
