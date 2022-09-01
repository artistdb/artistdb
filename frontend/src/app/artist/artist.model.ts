export interface ArtistInput {
    firstName: string;
    lastName: string;
    artistName?: string;
    pronouns?: string[];
    // make dateOfBirth?: [Date];
    placeOfBirth?: string;
    nationality?: string;
    language?: string;
    facebook?: string;
    instagram?: string;
    bandcamp?: string;
    bioGer?: string;
    bioEn?: string;
}

export const ARTIST_FIELDS = `{
    "data": [
        {
            "key": "First Name", 
            "required": true,
            "controlType": "textbox",
            "type": ""
        },
        {
            "key": "Last Name", 
            "required": true,
            "controlType": "textbox",
            "type": ""
        },
        {
            "key": "Artist Name", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        },
        {
            "key": "Pronouns", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        },
        {
            "key": "Place Of Birth", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        },
        {
            "key": "Nationality", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        },
        {
            "key": "Language", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        },
        {
            "key": "Facebook", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        },
        {
            "key": "Instagram", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        },
        {
            "key": "Bandcamp", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        },
        {
            "key": "Bio (German)", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        },
        {
            "key": "Bio (English)", 
            "required": false,
            "controlType": "textbox",
            "type": ""
        }
    ]
}`

export const MOCK_ARTISTS: ArtistInput[] = [
    {
        firstName: "Jens",
        lastName: "Rainer",
        artistName: "jens-rainer",
        pronouns: ["he", "him"],
        placeOfBirth: "drüben",
        nationality: "drüber",
        language: "zeug",
        facebook: "www.meta.com/bs",
        instagram: "www.meta.com/instant_bs",
        bandcamp: "www.bandcamp.com/zeug_xxx",
        bioGer: "Ich komm von drüben",
        bioEn: "I came from yonder",
    },
    {
        firstName: "Kai",
        lastName: "Uwe",
        artistName: "uwe-kai",
        pronouns: ["they", "them"],
        placeOfBirth: "drüben",
        nationality: "drüber",
        language: "zeug",
        facebook: "www.meta.com/bs",
        instagram: "www.meta.com/instant_bs",
        bandcamp: "www.bandcamp.com/zeug_xxx",
        bioGer: "Ich komm von drüben",
        bioEn: "I came from yonder",
    },
]