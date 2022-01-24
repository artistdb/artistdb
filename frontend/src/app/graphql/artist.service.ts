import { Injectable } from '@angular/core';
import { Mutation, gql } from 'apollo-angular';

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
