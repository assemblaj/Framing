# Framing
Framing is a value-key database that allows the user to load a JSON structure,
find values with respect to multiple contexts and traverse upward in their hierarchies.

## Frames
A 'Frame' is a set of contextual data about a specific value. This includes:     
**Subject** - The key for which this specific instance of the value was found  
**MetaData** - The keys that accompany the value's parent object, which can be used to access the other values in the object  
**data** - Values for the base object  
**Ancestors** - All ancestors of the base object  

## Example 
Example JSON: 
```
{
     "crops": [ 
       {
            "name": "apple", 
            "environment" : "orchard", 
            "grown_on": "trees", 
            "harvest_time" : 2.5
       }, 
       {
           "name": "orange", 
           "environment" : "farm", 
           "grown_on": "trees", 
           "harvest_time" : 1.5
       }
     ],

    "companies": [
        {
            "company": "Apple", 
            "NASDAQ": "AAPL",
            "Revenue" :	265.595,
            "Operating income" : 70.898,
            "Net income" : 9.531,
            "Total assets"	: -365.725,
            "Total equity"	: -107.147
        }, 
       {
        "company": "Microsoft", 
        "NASDAQ": "MSFT", 
        "Revenue" :	110.36,
        "Operating" : 35.05,
        "Net income" : 16.57,
        "Total assets"	: 58.84,
        "Total equity" : 82.71 
       }
    ], 

    "products": [
        {
            "flavor" : "candy apple", 
            "price" : 2.50, 
            "amount" : 2, 
            "discount" : 0, 
            "size"  : "normal"
        }, 
      {
        "flavor" : "chocolate", 
        "price" : 5.50, 
        "amount" : 1, 
        "discount" : 0, 
        "size"  : "medium"
      }
    ]    
}
```
