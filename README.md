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

GroupByMetaData("Apple") would return the Frames for all 
values that contain Apple (the company, the product, and the crop), 
as well as their parent objects, allowing the user to process 
siblings and ancestors. 

```
company|NASDAQ|Revenue|Operating income|Net income|Total assets|Total equity:[
	«Subject:company, MetaData:[company NASDAQ Revenue Operating income Net income Total assets Total equity]»
]
name|environment|grown_on|harvest_time:[
	«Subject:name, MetaData:[name environment grown_on harvest_time]»
]
flavor|price|amount|discount|size:[
	«Subject:flavor, MetaData:[flavor price amount discount size]»
]
```

## API 
### NewFramingDB 
```
Framing := NewFramingDB()
```

### DB.Load     
**input**: io.Reader    
**output** : error     

```
err = Framing.Load(r)
```

### DB.Get    
**input**: SearchValue struct, which currently supports    
 - value : string - Value to be searched   
 - exact : bool - Match only exact values or allow similar values?   
 - caseSentivite : bool - Match exact case or any case ?    

**output**:      
  - succes : bool - Whether the operation has succeeded or not      
  - frames : \*[]Frame - Pointer to list of frames      
```
exists, fs := Framing.Get([value])
fs = Slice of Frames 
```

### DB.GetDistinctMetaData   
**output**:
- value->frame map : map[string][]*Frame - 
```
fmap := Framing.GetDistincMetaData()
```

### DB.GetDistinct 
**input**:
- value : string - Value to be searched on  

**output**:
- succes : bool - Whether the operation has succeeded or not      
- frames : \*[]Frame - Pointer to list of frames that all have different MetaData (base object structure)    
```
e, fs := Framing.GetDistinct([value])
```

### DB.GroupByMetaData
**input**:
- valse : string - Value to be searched on 

**output**:
- succes : bool - Whether the operation has succeeded or not      
- metadata->frame map :  map[string][]\*Frame - Pointer to list of frames that all have different MetaData (base object structure)    

```
e, mdfmap := Framing.GroupByMetaData([value])
```

### Frame.Get
**input**: Key of base object (can be found in MetaData array)     
**output**: Value for that object as a string    
```
// Get values from the object 
exists, val = frame.Get([key])
```


## Potential Advancements 
- Multiple files 
- Multithreading 
- Network layer for inter-process communication,etc,  possibly via HTTP 
