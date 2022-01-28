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