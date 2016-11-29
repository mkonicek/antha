## Lesson 0: make mastermix and wire into PCR example

The protocol in this folder will generate instructions to set up PCR reactions, first making a mastermix for the PCR reactions and then using the mastermix for each reaction. You can run it by typing antharun
according to the "Reactiontoprimerpair" and "Reactiontotemplate" in the parameters file. 


To run an Antha Workflow with Antharun two types of files are needed: 


1. A worflow definition file
2. A parameters file 

### Workflow:
The workflow file specifies a set of Processes which call Antha elements 
(components) which are to be run. 


(B) Two elements wired together run once
 
![Auto PCR](autopcrhover.png)



### Parameters:
The parameters file assigns parameters for each of the processes specified in the workflow file

i.e. the parameters file is used to set the values for the input parameters.

The example parameters files in these folders show how to set variables specified in the parameters file to the actual values we want to assign to them.
One of the key variables you'll likely want to set are the liquid handling components (wtype.LHComponent) 


