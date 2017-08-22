from __future__ import division
import yaml
import sys
from collections import OrderedDict
import configparser
import re 
import os

### Data directory / volume mount point
datadir = '/datadir' 

### Workload template 
template = sys.argv[1]

### Config output folder
output = sys.argv[2]

### Fio job file name
paramfile = output + '/' + template.split("/")[-1].strip('yml') + 'vd'

### Load the params in input.YAML into a dict
with open(template, 'r') as stream:
    try:
        input_params = yaml.load(stream)
    except yaml.YAMLError as exc:
        print (exc)

### Define the actual vdbench params dict
vdbench_params = OrderedDict()

### Convert input params into vdbench params

## fsd params 

vdbench_params['fsd_name'] = "fsd-"+input_params['NAME']
vdbench_params['anchor'] = "/datadir"

# TODO How to fit into generic yaml, set to "1" by default 
vdbench_params['depth'] = '1' 
vdbench_params['width'] = '1'

vdbench_params['files'] = str(input_params['NUM_FILES'])
vdbench_params['size'] = str(input_params['FILE_SIZE']) + 'M'

if not input_params['BUFFERED_IO']:
    vdbench_params['openflags'] = 'o_direct'
else:
    vdbench_params['openflags'] = 'o_sync'

## fwd params

vdbench_params['fwd_name'] = "fwd-"+input_params['NAME']
vdbench_params['rdpct'] = str(input_params['RW_RATIO'])
vdbench_params['xfersize'] = str(input_params['IO_SIZE'])
vdbench_params['fileio'] = input_params['ACCESS_PATTERN'].lower()

# TODO Capping no of threads to file count if greater, need to notify user
if input_params['NUM_WORKERS'] <= input_params['NUM_FILES']:
    vdbench_params['threads'] = str(input_params['NUM_WORKERS'])
else: 
    vdbench_params['threads'] = str(input_params['NUM_FILES'])

vdbench_params['fileselect'] = input_params['ACCESS_PATTERN'].lower()

## rd params

vdbench_params['rd_name'] = "rd-"+input_params['NAME']
vdbench_params['elapsed'] = str(input_params['DURATION'])
vdbench_params['interval'] = "1"
vdbench_params['fwdrate'] = "max"
vdbench_params['format'] = "yes"
vdbench_params['warmup'] = str(input_params['WARMUP'])


### Construct the vdbench filesystem storage definition (fsd)

storage_definition= "fsd="+vdbench_params['fsd_name'] + "," + \
                    "anchor="+vdbench_params['anchor'] + ","+ \
                    "depth="+vdbench_params['depth'] + "," + \
                    "width="+vdbench_params['width'] + "," + \
                    "files="+vdbench_params['files'] + "," + \
                    "size="+vdbench_params['size']+ "," + \
                    "openflags="+vdbench_params['openflags'] + "\n" 



### Construct the vdbench filesystem workload definition (fwd)

workload_definition = "fwd="+vdbench_params['fwd_name'] + "," + \
                      "fsd="+vdbench_params['fsd_name'] + "," + \
                      "rdpct="+vdbench_params['rdpct'] + "," + \
                      "xfersize="+vdbench_params['xfersize'] + "," + \
                      "fileio="+vdbench_params['fileio'] + "," + \
                      "threads="+vdbench_params['threads'] + "," + \
                      "fileselect="+vdbench_params['fileselect'] + "\n"

#print "samples : %s " %(vdbench_params['elapsed'])

### Construct the vdbench run definition (rd) 

run_definition = "rd="+vdbench_params['rd_name'] + "," + \
                 "fwd="+vdbench_params['fwd_name'] + "," + \
                 "elapsed="+vdbench_params['elapsed'] + "," + \
                 "interval="+vdbench_params['interval'] + "," + \
                 "fwdrate="+vdbench_params['fwdrate'] + "," + \
                 "format="+vdbench_params['format'] + "," + \
                 "warmup="+vdbench_params['warmup'] + "\n"


### Construct the vdbench general definition (rd) after converting I/Ps

## gd params

if input_params['DEDUPE_PERCENTAGE'] != 0: 

    vdbench_params['dedupratio'] = "{:.1f}".format(100/(100-int(input_params['DEDUPE_PERCENTAGE'])))
    vdbench_params['dedupunit'] = str(input_params['IO_SIZE'])

    if input_params['COMPRESS_PERCENTAGE'] != 0:
        vdbench_params['compratio'] = "{:.1f}".format(100/(100-int(input_params['COMPRESS_PERCENTAGE'])))
    else:
        vdbench_params['compratio'] = '1'
 
    general_definition = "dedupratio="+vdbench_params['dedupratio'] + "," + \
                         "dedupunit="+vdbench_params['dedupunit'] + "," + \
                         "compratio="+vdbench_params['compratio'] + "\n"

else:
    if input_params['COMPRESS_PERCENTAGE'] != '0':
        vdbench_params['compratio'] = "{:.1f}".format(100/(100-int(input_params['COMPRESS_PERCENTAGE'])))
        general_definition = "compratio="+vdbench_params['compratio'] + "\n"
        
### Create the .vd paramfile         

fd = open(paramfile, 'w+')
if general_definition:
    fd.write(general_definition)   
fd.write(storage_definition)
fd.write(workload_definition)
fd.write(run_definition)
fd.close()









