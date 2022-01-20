import { Injectable } from '@angular/core';
import { Mutation, gql } from 'apollo-angular';

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

@Injectable({
  providedIn: 'root'
})
export class UpsertArtists extends Mutation {
  override document = gql`
    mutation upsertArtists($artists: [ArtistInput!]) {
      upsertArtists(input: $artists) {
        id
        firstName
        lastName
      }
    }
  `;
}
