Tests of mix commands
---------------------

Part 1: destination autoallocation - do we mix in place or not

This matrix defines how the three instruction types should behave given the 
parameters specified here when it comes to assigning destinations

the result depends on which instruction type we use, whether c1 is a sample or not
and whether a destination well is specified

instruction	c1 sample? 	dest well? 	RESULT

// Mix and MixInto will not move component 1 if it is not a sample

Mix		NO		N/A		Mix in place: move c2 onto c1
Mix		YES		N/A		move both c1 and c2 to a new, autoallocated plate and well 

MixInto		NO		NO		If C1 not in plate specified, error, otherwise move c2 on top
MixInto		NO		YES		If C1 not in plate and well specified, error, otherwise move c2 on top
MixInto		YES		NO		C1, C2 moved to plate specified, autoallocated well 
MixInto		YES		YES		C1, C2 moved to plate specified at given well location

// MixNamed always moves the two components

MixNamed	NO		NO		C1 + C2 moved to new plate with name prefix specified
MixNamed	NO		YES		C1 + C2 moved to new plate with name prefix, well specified
MixNamed	YES		NO		C1 + C2 moved to new plate with name prefix specified
MixNamed	YES		YES		C1 + C2 moved to new plate with name prefix, well specified
