job "demo-re" {
	datacenters = ["dc1"]

	# Restrict our job to only linux. We can specify multiple
	# constraints as needed.
	constraint {
		attribute = "${attr.kernel.name}"
		value = "linux"
	}

	group "test" {
		# Define the controller task to run
		task "do-nothing" {
			# Use a docker wrapper to run the task.
			driver = "raw_exec"

			env {
				PAUSE_TIME = "10"
			}

			config {
				command = "/vagrant/scripts/launch_re"
			}

			# We must specify the resources required for
			# this task to ensure it runs on a machine with
			# enough capacity.
			resources {
				cpu = 500 # 500 MHz
				memory = 256 # 256MB
				network {
					mbits = 20
				}
			}

		}
	}
}
