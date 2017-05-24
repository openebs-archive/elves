# There can only be a single job definition per file.
# Create a job with ID and Name 'example'
job "demo-vsm" {
	# Run the job in the global region, which is the default.
	# region = "global"

	# Specify the datacenters within the region this job can run in.
	datacenters = ["dc1"]

	# Restrict our job to only linux. We can specify multiple
	# constraints as needed.
	constraint {
		attribute = "${attr.kernel.name}"
		value = "linux"
	}

	# Create a 'jiva-ctl' group. Each task in the group will be
	# scheduled onto the same machine.
	group "jiva-ctl" {
		# Control the number of instances of this group.
		# Defaults to 1
		# count = 1

		# Configure the restart policy for the task group. If not provided, a
		# default is used based on the job type.
		restart {
			mode = "fail"
		}

		# Define a task to run
		task "ctl" {
			# Use Docker to run the task.
			driver = "docker"

			# Configure Docker driver with the image
			config {
				image = "openebs/jiva:latest"
                                network_mode = "host"
                                privileged = true
				command = "launch"
				args = [ 
                                         "controller", 
                                         "--frontend", "gotgt", 
                                         "--frontendIP", "172.28.128.101", 
                                         "vol1", "1g" 
                                       ]
				port_map = {
					iscsi = 3260
					api = 9501
				}
				logging {
					type = "journald"
				}
			}

			# We must specify the resources required for
			# this task to ensure it runs on a machine with
			# enough capacity.
			resources {
				cpu = 500 # 500 MHz
				memory = 256 # 256MB
				network {
					mbits = 20
					port "iscsi" {}
					port "api" {}
				}
			}
		}
	}

	# Create a 'jiva-rep' group. Each task in the group will be
	# scheduled onto the same machine.
	group "jiva-rep" {
		# Control the number of instances of this group.
		# Defaults to 1
		# count = 1

		# Configure the restart policy for the task group. If not provided, a
		# default is used based on the job type.
		restart {
			mode = "fail"
		}

		# Define a task to run
		task "rep1" {
			# Use Docker to run the task.
			driver = "docker"

			# Configure Docker driver with the image
			config {
				image = "openebs/jiva:latest"
                                network_mode = "host"
				command = "launch"
				args = [ 
                                         "replica", 
                                         "--frontendIP", "172.28.128.101", 
                                         "--listen", "172.28.128.102:9502", 
                                         "--size", "1g", 
                                         "vol1"
                                       ]
				port_map = {
					api = 9502
					res1 = 9503
					res2 = 9504
				}
				volumes = [
					"/tmp/jiva/rep1:/vol1"
				]
				logging {
					type = "journald"
				}
			}

			# We must specify the resources required for
			# this task to ensure it runs on a machine with
			# enough capacity.
			resources {
				cpu = 500 # 500 MHz
				memory = 256 # 256MB
				network {
					mbits = 20
					port "api" {}
					port "res1" {}
					port "res2" {}
				}
			}
		}
	}
}
