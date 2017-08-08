
A sonar device is installed in an submarine (vessel, aka kubernetes cluster) that detects submarines (aka custom-resources) that come in its vicinity. A sonar operator (aka custom-controller) job is to keep looking at the sonar (monitor), get the details of the detected submarines and take appropriate actions. For the sake of this example, let us assume the sonar-operator is trying to raise an alert (aka. publish event), whenever a non-US submarine is detected. 

