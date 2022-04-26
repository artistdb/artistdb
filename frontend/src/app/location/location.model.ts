export interface Location {
    name: String,
    country?: String,
    zip?: String,
    city?: String,
    street?: String,
    picture?: String,
    description?: String,
    lat?: String,
    lon?: String,
}

export const LOCATION_FIELDS = `{
    "data": [
        {
            "key": "name", 
            "required": true,
            "controlType": "textbox",
            "type": ""
        }, 
        {
            "key": "country",
            "required": false,
            "controlType": "textbox",
            "type": ""
        },
        {
            "key": "zip", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        }, 
        {
            "key": "city", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        }, 
        {
            "key": "street", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        }, 
        {
            "key": "picture", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        }, 
        {
            "key": "description", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        }, 
        {
            "key": "lat", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        }, 
        {
            "key": "lon", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        }
    ]
}`