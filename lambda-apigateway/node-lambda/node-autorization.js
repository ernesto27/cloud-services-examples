export const handler = async(event) => {
    console.log("Version 0.0.1");
    console.log("event: ", event)
     let response = {
         "isAuthorized": false,
         "context": {
             "stringKey": "value",
             "numberKey": 1,
             "booleanKey": true,
             "arrayKey": ["value1", "value2"],
             "mapKey": {"value1": "value2"}
         }
     };
     
     if (event.headers.authorization === "secretToken") {
         console.log("allowed");
         response = {
             "isAuthorized": true,
             "context": {
                 "stringKey": "value",
                 "numberKey": 1,
                 "booleanKey": true,
                 "arrayKey": ["value1", "value2"],
                 "mapKey": {"value1": "value2"}
             }
         };
     }
 
     return response;
 
 };
 