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
jobfile = output + '/' + template.split("/")[-1].strip('yml') + 'fio'

### Define config parser object
config = configparser.ConfigParser()

### Load the params in input.YAML into a dict
with open(template, 'r') as stream:
    try:
        input_params = yaml.load(stream)
    except yaml.YAMLError as exc:
        print (exc)

### Define the actual fio params dict
fio_params = OrderedDict()

### Add directly replaceable param key:value pairs from input.yaml into fio_params dict
fio_params['name'] = input_params['NAME']
fio_params['directory'] = datadir
fio_params['bs'] = input_params['IO_SIZE']
fio_params['rwmixread'] = input_params['RW_RATIO']
fio_params['numjobs'] = input_params['NUM_WORKERS']
fio_params['nrfiles'] = input_params['NUM_FILES']
fio_params['filesize'] = str(input_params['FILE_SIZE']) + 'M'
fio_params['ramp_time'] = input_params['WARMUP']
fio_params['runtime'] = input_params['DURATION']
fio_params['blockalign'] = input_params['IO_ALIGNMENT']
fio_params['dedupe_percentage'] = input_params['DEDUPE_PERCENTAGE']
fio_params['buffer_compress_percentage'] = input_params['COMPRESS_PERCENTAGE']

### Add conditional params from input.yaml into fio_params dict 

# Determine I/O access pattern 
if input_params['ACCESS_PATTERN'] == 'Random':

    if input_params['RW_RATIO'] == '0':
        fio_params['readwrite'] = 'randwrite'

    elif input_params['RW_RATIO'] == '100':
        fio_params['readwrite'] = 'randread'    

    else:
        fio_params['readwrite'] = 'randrw'
    
    fio_params['randrepeat'] = '0'
    fio_params['norandommap'] = '1'
    fio_params['refill_buffers'] = '1'

elif input_params['ACCESS_PATTERN'] == 'Sequential':
    
    if input_params['RW_RATIO'] == '0':
        fio_params['readwrite'] = 'write'
    
    elif input_params['RW_RATIO'] == '100':
        fio_params['readwrite'] = 'read'

    else:
        fio_params['readwrite'] = 'rw'

# Determine I/O fs cache usage
if not input_params['BUFFERED_IO']:
    fio_params['buffered'] = '0'
    fio_params['invalidate'] = '1'

# Determine Data transfer type 
if input_params['DATA_TRANSFER'] == "async":
    fio_params['ioengine'] = 'libaio'
    fio_params['iodepth'] = input_params['QUEUE_DEPTH']

elif input_params['DATA_TRANSFER'] == "sync":
    fio_params['ioengine'] = 'sync'
#   fio_params['fsync'] = '1'

# Determine I/O arrival pattern
if input_params['BURST_IO']:
    fio_params['rate_process'] = 'poisson'

### Add fio execution defaults (run type, reporting & logging) into fio_params dict
fio_params['time_based'] = '1'
fio_params['group_reporting'] = '1'
fio_params['per_job_logs'] = '0' 
fio_params['write_iops_log'] = 'iops.log'
fio_params['write_bw_log'] = 'bw.log'
fio_params['write_lat_log'] = 'lat.log'

### Write the fio_params into test.fio job file
config['job'] = fio_params
f = open(jobfile, 'w')
config.write(f)
f.close()

### Sanitize the fio config file
file_handle = open(jobfile, 'rb')
file_string = file_handle.read()
file_handle.close()
file_string = (re.sub(" = ", "=", file_string))
file_handle = open(jobfile, 'wb')
file_handle.write(file_string)
file_handle.close()






