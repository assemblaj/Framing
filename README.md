# Framing
Framing is a value-key database that allows the user to load a JSON structure, 
find values with respect to multiple contexts and traverse upward in their hierarchies.  

## Frames 
A 'Frame' is a set of contextual data about a specific value. This includes:

Subject - The key for which this specific instance of the value was found
MetaData - The  keys that accompany the value's parent object, which can 
           be used to access the other values in the object
data - Values for the base object 
Ancestors - All ancestors of the base object 
