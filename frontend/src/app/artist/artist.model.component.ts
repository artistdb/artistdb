export interface ArtistInput {
    firstName: string;
    lastName: string;
    artistName?: string;
    pronouns?: [string];
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
    },
    {
        firstName: "Kai",
        lastName: "Uwe",
        artistName: "uwe-kai"
    },
]