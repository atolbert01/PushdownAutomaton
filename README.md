# RESTful PDA Processor
Author: Aaron Tolbert-Smith
Date: 04/12/2020
Class: CMSC 621 6193 (Spring 2020)


To run the server execute the following command inside the proj1 directory:

	./proj1

You should receive the output:

	Listening on port 8080...

Then, navigate to the driver directory and open a different terminal. Once there, execute the
following command to run the tests:

	./driver

This will send a series of requests to the server. The output from both the server and driver are
saved in the file 'standard-output.txt' in the proj1 directory.

I had intended to include a graceful shutdown method for the server to be executed once the tests
are completed, but unfortunately ran out of time. So to stop the server, you will have to kill the
process. In the terminal from which ./proj1 was launched, press:

	Ctrl + C