import { TestBed } from '@angular/core/testing';

import { UpsertArtists } from './artist.service';

describe('UpsertArtists', () => {
  let service: UpsertArtists;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(UpsertArtists);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
