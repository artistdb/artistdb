import { TestBed } from '@angular/core/testing';

import { FieldControlService } from './field-control.service';

describe('FieldControlService', () => {
  let service: FieldControlService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(FieldControlService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
